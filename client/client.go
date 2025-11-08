package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func StartClient(addr string) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+addr+"/ws", nil)
	if err != nil {
		log.Fatal("Failed to connect.", err)
	}
	defer conn.Close()

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

	input := bufio.NewScanner(os.Stdin)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			fmt.Println("\nDisconnected.")
			return

		default:
			if input.Scan() {
				text := input.Text()
				if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
					return
				}
			}
		}
	}

}
