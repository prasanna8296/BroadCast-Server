package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader converts an incoming HTTP request into a WebSocket connection.
// CheckOrigin always returns true so that clients from any origin may connect.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket upgrades the HTTP request to a WebSocket connection, registers
// the new client with the Hub, and reads messages from the client in a goroutine,
// forwarding each one to the Hub for broadcast to all connected clients.
func HandleWebSocket(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// Register the connection so it starts receiving broadcasts.
	hub.register <- conn

	// Read loop: runs until the client disconnects or an error occurs.
	go func() {
		defer func() {
			// Ensure the client is removed from the Hub on exit.
			hub.unregister <- conn
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break

			}
			// Forward the received message to the Hub for broadcasting.
			hub.broadcast <- msg
		}
	}()
}
