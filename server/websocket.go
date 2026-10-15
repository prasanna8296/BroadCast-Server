package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// writeWait bounds how long a single write to a client may block, so one
	// unresponsive client cannot stall the Hub's broadcast loop indefinitely.
	writeWait = 5 * time.Second
	// pongWait is how long we wait for a client pong before treating the
	// connection as dead. It must be larger than pingPeriod.
	pongWait = 60 * time.Second
	// pingPeriod is how often the server pings a client to keep the connection
	// alive and to detect half-open connections. It must be smaller than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// upgrader converts an incoming HTTP request into a WebSocket connection.
// CheckOrigin always returns true so that clients from any origin may connect.
// Note: this is permissive and fine for a local demo, but in production you
// should validate the request Origin to prevent cross-site WebSocket hijacking.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// pingClient keeps a client connection alive and detects half-open connections
// by sending periodic pings until the write fails or the connection closes via
// writerDone. It runs alongside the Hub's broadcast writes to the same conn; the
// WriteControl call used by WriteMessage is safe for concurrent use.
func pingClient(conn *websocket.Conn, writerDone <-chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-writerDone:
			return
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
				return
			}
		}
	}
}

// HandleWebSocket upgrades the HTTP request to a WebSocket connection, registers
// the new client with the Hub, and reads messages from the client in a goroutine,
// forwarding each one to the Hub for broadcast to all connected clients.
func HandleWebSocket(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Upgrade already wrote an error response, so there is nothing to do:
		// the connection was never established and must not touch the Hub.
		return
	}
	// Register the connection so it starts receiving broadcasts.
	hub.register <- conn

	// writerDone is closed when the read loop ends, signalling the pinger to stop.
	writerDone := make(chan struct{})

	// Read loop: runs until the client disconnects or an error occurs.
	go func() {
		defer func() {
			// Ensure the client is removed from the Hub on exit.
			hub.unregister <- conn
			close(writerDone)
		}()

		// If no pong arrives within pongWait, ReadMessage will fail and the
		// loop will exit, letting the Hub drop the dead connection.
		conn.SetReadDeadline(time.Now().Add(pongWait))

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// Each successful read proves the client is alive: extend the deadline.
			conn.SetReadDeadline(time.Now().Add(pongWait))
			// Forward the received message to the Hub for broadcasting.
			hub.broadcast <- msg
		}
	}()

	// Keep-alive loop: pings the client so idle connections are not dropped.
	go pingClient(conn, writerDone)
}
