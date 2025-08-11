package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"websocket-service/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	repo storage.MessageRepository
	ID   string
}

func ServeWS(hub *Hub, repo storage.MessageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// получаем идентификатор (открытый ключ) из query
		id := c.Query("id")
		if id == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade:", err)
			return
		}

		client := &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
			repo: repo,
			ID:   id,
		}
		hub.register <- client

		go client.writePump()
		client.readPump()
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg struct {
			To       string    `json:"to"`
			From     string    `json:"from"`
			Data     []byte    `json:"data"`
			SendTime time.Time `json:"send_time"`
		}
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("Invalid msg:", err)
			continue
		}
		c.hub.send <- Envelope{To: msg.To, From: msg.From, Message: msg.Data, SendTime: time.Now()}
		//go func() {
		//	message := models.Message{
		//		MessageId:      uuid.New(),
		//		SenderId:       msg.From,
		//		ConversationId: msg.To,
		//		Ciphertext:     string(msg.Data),
		//		Nonce:          "",
		//		CreatedAt:      time.Now(),
		//		IsRead:         false,
		//		IsDeleted:      false,
		//	}
		//	ctx := context.Background()
		//	err = c.repo.InsertMSG(ctx, message)
		//	if err != nil {
		//		log.Println("Ошибка добавления сообщения в базу:", err)
		//	}
		//}()
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msgBytes, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				log.Println("Error sending msg:", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
