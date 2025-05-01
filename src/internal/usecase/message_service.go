package services

import (
	"context"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	usecase_errors "libs/src/internal/usecase/errors"
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

func (s *MessageService) SendMessage(ctx context.Context, sender dto.UserDTO, messageRequest dto.SendMessageRequest, chatId int64) (*dto.MessagePreviewDTO, error) {
	if sender.Role == enums.ANONYMOUS || !sender.IsActive {
		return &dto.MessagePreviewDTO{}, usecase_errors.UnauthorizedError{Msg: "You must be logged in to send a message"}
	}

	members, err := s.ChatMemberRepository.Filter(ctx, "chat_id = ? AND user_id = ?", chatId, sender.ID)
	if err != nil {
		return &dto.MessagePreviewDTO{}, err
	}
	if len(members) != 1 {
		return &dto.MessagePreviewDTO{}, usecase_errors.BadRequestError{Msg: "You are not a member of this chat"}
	}

	if members[0].MemberRole < enums.MEMBER {
		return &dto.MessagePreviewDTO{}, usecase_errors.BadRequestError{Msg: "You cannot send a message to this chat"}
	}

	message := domain.NewMessageObject(sender.ID, chatId, messageRequest.Message)
	err = s.MessageRepository.Create(ctx, message)
	if err != nil {
		return &dto.MessagePreviewDTO{}, err
	}

	messagePreview := &dto.MessagePreviewDTO{
		Id:             message.Id.Hex(),
		Content:        message.Content,
		SenderUsername: sender.Username,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
	}
	return messagePreview, nil
}
