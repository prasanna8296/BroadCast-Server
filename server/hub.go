package server

import (
    "github.com/gorilla/websocket"
)

// Hub maintains the set of active WebSocket clients and broadcasts
// messages received from any client to all connected clients.
type Hub struct {
	clients    map[*websocket.Conn]bool // registered clients
	broadcast  chan []byte              // inbound messages from clients
	register   chan *websocket.Conn     // requests to add a client
	unregister chan *websocket.Conn     // requests to remove a client
}

// NewHub creates and returns a Hub with its channels and client set initialized.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run starts the Hub's event loop. It must be started in its own goroutine
// before clients connect. The loop serializes access to the clients map,
// handling registrations, unregistrations, and broadcasts.
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			// A new client connected; add it to the set.
			h.clients[conn] = true
		case conn := <-h.unregister:
			// A client disconnected; remove and close it.
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
		case message := <-h.broadcast:
			// Relay the message to every connected client.
			for conn := range h.clients {

				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					// Write failed, so the client is gone; drop it.
					delete(h.clients, conn)
					conn.Close()
				}
			}
		}
	}
}