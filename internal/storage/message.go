package storage

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"websocket-service/internal/models"
)

type MessageRepository interface {
	InsertMSG(ctx context.Context, msg models.Message) error
}

type repoMSG struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &repoMSG{db: db}
}

func (r *repoMSG) InsertMSG(ctx context.Context, msg models.Message) error {
	err := r.db.WithContext(ctx).Save(&msg).Error
	fmt.Println(err)
	return err
}
