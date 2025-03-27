package services

import (
	"libs/src/internal/repositories"
	"libs/src/settings"
)

type ChatService struct {
	UserRepository       *repositories.UserRepository
	ChatRepository       *repositories.ChatRepository
	ChatMemberRepository *repositories.ChatMemberRepository
}

func NewChatService(app *settings.App) *ChatService {
	return &ChatService{
		UserRepository:       repositories.NewUserRepository(app),
		ChatRepository:       repositories.NewChatRepository(app),
		ChatMemberRepository: repositories.NewChatMemberRepository(app),
	}
}

//func (s *ChatService) CreateChat() ([]domain.Chat, error) {
//	return s.ChatRepository.GetChatsByUserID(userID)
//}
