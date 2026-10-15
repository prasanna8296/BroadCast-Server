package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// StartServer creates a Hub, starts its event loop in a goroutine, and starts
// an HTTP server on the given port that upgrades /ws requests to WebSockets.
func StartServer(port int) {

	// The Hub runs in its own goroutine so it can service register,
	// unregister, and broadcast events independently of incoming requests.
	hub := NewHub()
	go hub.Run()

	// Gin router; the only route is the WebSocket endpoint.
	// gin.Default() includes the logger and recovery middleware, which is
	// convenient for development.
	r := gin.Default()
	r.GET("/ws", func(ctx *gin.Context) {
		HandleWebSocket(hub, ctx)
	})

	addr := fmt.Sprintf(":%d", port)
	log.Println("Broadcast server running on ", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server error:", err)
	}

}
