package services

import (
	"libs/src/internal/repositories"
	"libs/src/settings"
)

type MessageService struct {
	App                  *settings.App
	MessageRepository    repositories.IMessageRepository
	ChatRepository       repositories.IChatRepository
	ChatMemberRepository repositories.IChatMemberRepository
}

func NewMessageService(app *settings.App) *MessageService {
	return &MessageService{
		App:                  app,
		MessageRepository:    repositories.NewMessageRepository(app),
		ChatRepository:       repositories.NewChatRepository(app),
		ChatMemberRepository: repositories.NewChatMemberRepository(app),
	}
}
