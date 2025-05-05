package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	BaseMongo
	SenderId  int64  `bson:"sender_id" json:"sender_id"`
	ChatId    int64  `bson:"chat_id" json:"chat_id"`
	Content   string `bson:"content" json:"content"`
	IsRead    bool   `bson:"is_read" json:"is_read"`
	IsUpdated bool   `bson:"is_updated" json:"is_updated"`
	IsDeleted bool   `bson:"is_deleted" json:"is_deleted"`
}

func NewMessageObject(senderId, chatId int64, content string) *Message {
	return &Message{
		BaseMongo: BaseMongo{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		SenderId: senderId,
		ChatId:   chatId,
		Content:  content,
	}
}
