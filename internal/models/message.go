package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	MessageId      uuid.UUID `json:"message_id"`
	SenderId       string    `json:"sender_id"`
	ConversationId string    `json:"conversation_id"`
	Ciphertext     string    `json:"ciphertext"`
	Nonce          string    `json:"nonce"`
	CreatedAt      time.Time `json:"created_at"`
	IsRead         bool      `json:"is_read"`
	IsDeleted      bool      `json:"is_deleted"`
}

func (_ *Message) TableName() string {
	return "messages"
}
