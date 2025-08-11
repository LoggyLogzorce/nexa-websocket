package main

import (
	"log"
	"websocket-service/internal/db"
	"websocket-service/internal/storage"
	"websocket-service/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	dbConn := db.Connect()
	repoMSG := storage.NewMessageRepository(dbConn)

	r := gin.Default()
	hub := ws.NewHub()
	go hub.Run()

	r.GET("/ws", ws.ServeWS(hub, repoMSG))

	log.Println("Server running on :8080")
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
