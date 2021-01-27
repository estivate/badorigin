package badorigin

import (
	"net/http"
	"strings"
)

// NoDirFS creates a handler for a File Server that doesn't
// return directory listings.
func NoDirFS(root_dir, path string) http.HandlerFunc {
	fs := http.FileServer(CustomFS{http.Dir(root_dir)})
	return http.StripPrefix(strings.TrimRight(path, "/"), fs).ServeHTTP
}

type CustomFS struct {
	fs http.FileSystem
}

func (fs CustomFS) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
