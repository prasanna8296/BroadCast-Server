package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

// StartClient connects to a broadcast server at the given address (host:port),
// prints messages received from the server, and sends lines typed on stdin to
// the server until the connection closes or the user presses Ctrl+C.
func StartClient(addr string) {
	// Dial the server's WebSocket endpoint.
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+addr+"/ws", nil)
	if err != nil {
		log.Fatal("Failed to connect.", err)
	}
	defer conn.Close()

	// Receiver goroutine: prints every broadcast message; closes done on exit.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			fmt.Printf(">> %s\n", message)
		}
	}()
	fmt.Println("Connected to the broadcast server")
	fmt.Println("Type messages and press Enter to send.")
	fmt.Println("Press Ctrl+C to exit.")

	// Set up stdin reading and Ctrl+C handling.
	input := bufio.NewScanner(os.Stdin)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-done:
			// The server closed the connection.
			return

		case <-interrupt:
			// The user asked to quit.
			fmt.Println("\nDisconnected.")
			return

		default:
			// Read a line from stdin and send it to the server.
			if input.Scan() {
				text := input.Text()
				if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
					return
				}
			}
		}
	}

}
