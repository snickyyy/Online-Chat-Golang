package services

import (
	"context"
	"errors"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
)

type ChatService struct {
	App                  *settings.App
	UserRepository       repositories.IUserRepository
	ChatRepository       repositories.IChatRepository
	ChatMemberRepository repositories.IChatMemberRepository
}

func NewChatService(app *settings.App) *ChatService {
	return &ChatService{
		App:                  app,
		UserRepository:       repositories.NewUserRepository(app),
		ChatRepository:       repositories.NewChatRepository(app),
		ChatMemberRepository: repositories.NewChatMemberRepository(app),
	}
}

func (s *ChatService) CreateChat(ctx context.Context, request dto.CreateChatRequest, user dto.UserDTO) (dto.ChatDTO, error) {
	if user.Role == enums.ANONYMOUS || !user.IsActive {
		return dto.ChatDTO{}, usecase_errors.UnauthorizedError{Msg: "You must be logged in to create members"}
	}

	newChat := domain.Chat{
		Title:       request.Title,
		Description: request.Description,
		OwnerID:     user.ID,
	}

	err := s.ChatRepository.Create(ctx, &newChat)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return dto.ChatDTO{}, usecase_errors.AlreadyExistsError{Msg: "Chat with this name already exists"}
		}
		return dto.ChatDTO{}, err
	}
	return newChat.ToDTO(), nil
}

func (s *ChatService) DeleteChat(ctx context.Context, caller dto.UserDTO, chatID int64) error {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return usecase_errors.UnauthorizedError{Msg: "You must be logged in to delete members"}
	}

	chat, err := s.ChatRepository.GetById(ctx, chatID)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.NotFoundError{Msg: "Chat with this ID not found"}
		}
		return err
	}

	if chat.OwnerID != caller.ID {
		return usecase_errors.PermissionError{Msg: "You have no permission to delete this chat"}
	}
	err = s.ChatRepository.DeleteById(ctx, chatID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return usecase_errors.NotFoundError{Msg: "Chat with this ID not found"}
		}
		return err
	}
	return nil
}

func (s *ChatService) ChangeChat(ctx context.Context, caller dto.UserDTO, chatId int64, request dto.ChangeChatRequest) error {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return usecase_errors.UnauthorizedError{Msg: "You must be logged in to change members"}
	}

	chat, err := s.ChatRepository.GetById(ctx, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.NotFoundError{Msg: "Chat with this ID not found"}
		}
		return err
	}
	if chat.OwnerID != caller.ID {
		return usecase_errors.PermissionError{Msg: "You have no permission to change this chat"}
	}

	filterData := map[string]*string{
		"title":       request.NewTitle,
		"description": request.NewDescription,
	}

	updateData := make(map[string]any, len(filterData))

	for k, v := range filterData {
		if v != nil {
			updateData[k] = v
		}
	}

	err = s.ChatRepository.UpdateById(ctx, chatId, updateData)
	if err != nil {
		if errors.As(err, &repositories.ErrDuplicate) {
			return usecase_errors.AlreadyExistsError{Msg: "Chat with this name already exists"}
		}
	}
	return nil
}

func (s *ChatService) GetListForUser(ctx context.Context, caller dto.UserDTO, page int) ([]dto.ChatDTO, error) {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return []dto.ChatDTO{}, nil
	}

	if page < 1 {
		return []dto.ChatDTO{}, usecase_errors.BadRequestError{Msg: "Invalid page"}
	}

	list, err := s.ChatRepository.GetListForUser(ctx, caller.ID, s.App.Config.Pagination.ChatList, (page-1)*s.App.Config.Pagination.ChatList)
	if err != nil {
		if errors.As(err, &repositories.ErrLimitMustBePositive) || errors.As(err, &repositories.ErrOffsetMustBePositive) {
			return []dto.ChatDTO{}, usecase_errors.BadRequestError{Msg: "Invalid page"}
		}
		return []dto.ChatDTO{}, err
	}

	result := make([]dto.ChatDTO, len(list))

	for i, v := range list {
		result[i] = v.ToDTO()
	}
	return result, nil
}

func (s *ChatService) Search(ctx context.Context, caller dto.UserDTO, name string, page int) ([]dto.ChatDTO, error) {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return []dto.ChatDTO{}, nil
	}

	if page < 1 {
		return []dto.ChatDTO{}, usecase_errors.BadRequestError{Msg: "Invalid page"}
	}

	list, err := s.ChatRepository.SearchForUser(ctx, caller.ID, name, s.App.Config.Pagination.ChatList, (page-1)*s.App.Config.Pagination.ChatList)
	if err != nil {
		if errors.As(err, &repositories.ErrLimitMustBePositive) || errors.As(err, &repositories.ErrOffsetMustBePositive) {
			return []dto.ChatDTO{}, usecase_errors.BadRequestError{Msg: "Invalid page"}
		}
		return []dto.ChatDTO{}, err
	}

	result := make([]dto.ChatDTO, len(list))

	for i, v := range list {
		result[i] = v.ToDTO()
	}
	return result, nil
}

func (s *ChatService) GetById(ctx context.Context, caller dto.UserDTO, chatId int64) (dto.ChatDTO, error) {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return dto.ChatDTO{}, usecase_errors.UnauthorizedError{Msg: "You must be logged in to get members"}
	}

	member, err := s.ChatMemberRepository.Filter(ctx, "user_id = ? AND chat_id = ?", caller.ID, chatId)
	if err != nil {
		return dto.ChatDTO{}, err
	}
	if len(member) < 1 {
		return dto.ChatDTO{}, usecase_errors.NotFoundError{Msg: "Chat with this ID not found"}
	}

	chat, err := s.ChatRepository.GetById(ctx, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return dto.ChatDTO{}, usecase_errors.NotFoundError{Msg: "Chat with this ID not found"}
		}
		return dto.ChatDTO{}, err
	}
	return chat.ToDTO(), nil
}
