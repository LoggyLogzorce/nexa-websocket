package main

import (
	"log"
	"websocket-service/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	hub := ws.NewHub()
	go hub.Run()

	// клиент подключается так: ws://0.0.0.0:8080/ws?id=<your_pubkey>
	r.GET("/ws", ws.ServeWS(hub))

	log.Println("Server running on :8080")
	r.Run("0.0.0.0:8080")
}
