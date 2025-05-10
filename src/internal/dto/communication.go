package dto

import "encoding/json"

type ChatCommunication struct {
	Action string          `json:"action" binding:"required"`
	Body   json.RawMessage `json:"Body"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type SendMessageRequest struct {
	ChatId  int64  `json:"chat_id"`
	Message string `json:"message" binding:"required, min=1,max=200"`
}

type SendMessage struct {
	ChatId      int64              `json:"chat_id" binding:"required"`
	MessageBody *MessagePreviewDTO `json:"body"`
}

type Subscribe struct {
	ChatId int64 `json:"chat_id" binding:"required"`
}
