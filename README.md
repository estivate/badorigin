
# Bad Origin ~ Flakey Test Servers for your ReverseProxy 

When working on Reverse Proxy projects it helps to have some "fake" origin 
servers to round robin content from. This library lets you easily spin up
servers that perform in different ways for development and testing purposes.

## Status

Spin up any number of servers on different ports, they'll all operate with
a random delay.

## Install

`go get -u github.com/estivate/bad-origin`

## Example


```go
package main

import (
	"flag"
	"fmt"

	"github.com/estivate/badorigin"
)

func main() {
	// are we testing? 
	var testMode bool
	flag.BoolVar(&testMode, "testMode", false, "run in test mode (default false)")
	flag.Parse()

	// if testing, launch the servers
	if testMode {
		bo := badorigin.NewServers(":8000", ":8001", ":8002")
		bo.LaunchServers()
	}

	// spin up your real work where you are running some kind of
    // server that blocks forever.
	fmt.Println("launching a reverse proxy in front of servers here")
	// block forever
	select {}
}

```

