package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer(port int) {

	hub := NewHub()
	go hub.Run()

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
