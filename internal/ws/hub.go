package ws

import (
	"encoding/json" // ← новый импорт
	"sync"
	"time"
)

// Envelope теперь хранит и кто шлёт
type Envelope struct {
	From     string
	To       string
	Message  []byte
	SendTime time.Time
}

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	send       chan Envelope
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		send:       make(chan Envelope),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c.ID] = c
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, c.ID)
			close(c.send)
			h.mu.Unlock()

		case env := <-h.send:
			h.mu.Lock()
			if dest, ok := h.clients[env.To]; ok {
				// Формируем JSON с полями from и data (data автоматически base64)
				out := struct {
					From string `json:"from"`
					Data []byte `json:"data"`
				}{
					From: env.From,
					Data: env.Message,
				}
				b, err := json.Marshal(out)
				if err == nil {
					select {
					case dest.send <- b:
					default:
						// при переполнении канала отключаем
						delete(h.clients, dest.ID)
						close(dest.send)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}
