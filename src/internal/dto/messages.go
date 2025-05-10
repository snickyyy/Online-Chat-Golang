package dto

import (
	"time"
)

type BaseMessageDTO struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	ChatId    int64     `json:"chat_id"`
	SenderId  string    `json:"sender_id"`
	IsEdited  bool      `json:"is_edited"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MessagePreviewDTO struct {
	Id             string    `json:"id"`
	Content        string    `json:"content"`
	SenderUsername string    `json:"sender_username"`
	IsEdited       bool      `json:"is_edited"`
	IsRead         bool      `json:"is_read"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}
