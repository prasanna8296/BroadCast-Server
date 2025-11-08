package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	hub.register <- conn

	go func() {
		defer func() {
			hub.unregister <- conn
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break

			}
			hub.broadcast <- msg
		}
	}()
}
