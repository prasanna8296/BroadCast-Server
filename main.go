package main

import (
	"broadcast/client"
	"broadcast/server"
	"fmt"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  broadcast-server start [port]")
		fmt.Println("  broadcast-server connect [host:port]")
		return
	}
	cmd := os.Args[1]

	switch cmd {
	case "start":
		port := 8080
		if len(os.Args) > 2 {
			if p, err := strconv.Atoi(os.Args[2]); err == nil {
				port = p
			}
		}
		server.StartServer(port)

	case "connect":
		addr := "localhost:8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		client.StartClient(addr)
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
