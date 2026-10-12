package main

import (
	"broadcast/client"
	"broadcast/server"
	"fmt"
	"os"
	"strconv"
)

// main is the entry point of the broadcast server CLI.
// It supports two commands:
//
//	start   [port]       - start the broadcast server, optionally on the given port
//	connect [host:port]  - connect to a running broadcast server as a client
func main() {

	// Without a command there is nothing to do, so print usage and exit.
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  broadcast-server start [port]")
		fmt.Println("  broadcast-server connect [host:port]")
		return
	}
	cmd := os.Args[1]

	switch cmd {
	case "start":
		// Default to port 8080, but allow an override via the second argument.
		port := 8080
		if len(os.Args) > 2 {
			if p, err := strconv.Atoi(os.Args[2]); err == nil {
				port = p
			}
		}
		server.StartServer(port)

	case "connect":
		// Default to a local server on port 8080, but allow an override.
		addr := "localhost:8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		client.StartClient(addr)
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
