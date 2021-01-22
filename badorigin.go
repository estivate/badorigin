// Package origin simulates upstream content servers
// for development and testing of reverse proxy functions
package badorigin

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Servers struct {
	Ports []string
	Debug bool
}

func NewServers(ports ...string) *Servers {
	fmt.Println("creating new servers")
	return &Servers{Ports: ports, Debug: false}
}

func (s *Servers) SetDebug() {
	fmt.Println("setting debug")
	s.Debug = true
}

//LaunchServers runs some downstream content servers for overdrive
// go grab stuff from.
func (s *Servers) LaunchServers() {
	fmt.Println("launching servers")

	// new gorilla mux router
	r := mux.NewRouter()

	changeHeaderThenServe := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Set some header.
			w.Header().Add("Server", "Test Origin")
			// Serve with the actual handler.
			h.ServeHTTP(w, r)
		}
	}

	fs := changeHeaderThenServe(http.StripPrefix("/sites/", http.FileServer(http.Dir("./testing/sites/"))))
	r.PathPrefix("/sites/").Handler(fs)

	// redirects
	r.PathPrefix("/redirect/{code}/{location}").HandlerFunc(redirectHandler)

	// notFound
	r.PathPrefix("/notfound").HandlerFunc(http.NotFound)

	// Errors
	r.PathPrefix("/error/{code}/{message}").HandlerFunc(errorHandler)

	// middleware
	r.Use(logging)
	r.Use(chaos)
	//r.Use(setCookies)

	http.Handle("/", r)

	// launch a few web servers
	for _, p := range s.Ports {
		if s.Debug {
			fmt.Printf("server spinning up on: %s\n", p)
		}
		port := p
		go func() { launchServer(port, r) }()
	}

	// block forever
	select {}
}

// middleware: chaos
// inserts response delays and such
func chaos(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		delay := rand.Intn(1000000)
		time.Sleep(time.Duration(delay) * time.Microsecond)

		f.ServeHTTP(w, r)

	})
}

// middleware: set a cookie
// let's set a cookie if they need it
func setCookies(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie := &http.Cookie{
			Name:  "bstest",
			Value: base64.StdEncoding.EncodeToString([]byte("/test home.html /whatever")),
		}

		http.SetCookie(w, cookie)

		f.ServeHTTP(w, r)

	})
}

// middleware: logging
// writes a simple log of what's going on
func logging(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		f.ServeHTTP(w, r)

		endTime := time.Now()
		diff := endTime.Sub(startTime)
		log.Printf("Origin: Returning %s took %v", r.Host+r.URL.Path, diff.Seconds())
	})
}

// handler - redirects
func redirectHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	code, _ := strconv.Atoi(vars["code"])
	location := fmt.Sprintf("http://%s/", vars["location"])

	http.Redirect(w, r, location, code)
}

// handler - errors
func errorHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	code, _ := strconv.Atoi(vars["code"])
	message := strings.Replace(vars["message"], "_", " ", -1)
	log.Printf("here: %s", message)
	http.Error(w, message, code)
}

func launchServer(port string, handler http.Handler) {
	if err := http.ListenAndServe(port, handler); err != nil {
		panic(err)
	}
}