package main

import (
	"flag"
	"fmt"

	"github.com/estivate/badorigin"
)

func main() {
	// use flags, configs, whatever to specify if you want to
	// launch these servers
	var testMode bool
	flag.BoolVar(&testMode, "testMode", false, "run in test mode (default false)")
	flag.Parse()

	// if testing, launch the servers
	if testMode {
		bo := badorigin.NewServers(":8000", ":8001", ":8002")
		bo.SetDebug()
		bo.LaunchServers()
	}

	//... spin up your real work
	fmt.Println("launching a reverse proxy in front of servers here")

	// block forever
	select {}
}
