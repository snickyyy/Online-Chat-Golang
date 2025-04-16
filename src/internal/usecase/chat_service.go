package services

import (
	"errors"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	api_errors "libs/src/internal/usecase/errors"
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

func (s *ChatService) CreateChat(request dto.CreateChatRequest, user dto.UserDTO) (dto.ChatDTO, error) {
	if user.Role == enums.ANONYMOUS || !user.IsActive {
		return dto.ChatDTO{}, api_errors.ErrUnauthorized
	}

	newChat := domain.Chat{
		Title:       request.Title,
		Description: request.Description,
		OwnerID:     user.ID,
	}

	err := s.ChatRepository.Create(&newChat)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return dto.ChatDTO{}, api_errors.ErrChatAlreadyExists
		}
		return dto.ChatDTO{}, err
	}
	return newChat.ToDTO(), nil
}

func (s *ChatService) DeleteChat(caller dto.UserDTO, chatID int64) error {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return api_errors.ErrUnauthorized
	}

	chat, err := s.ChatRepository.GetById(chatID)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return api_errors.ErrChatNotFound
		}
		return err
	}

	if chat.OwnerID != caller.ID {
		return api_errors.ErrNotEnoughPermissionsForDelete
	}
	err = s.ChatRepository.DeleteById(chatID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return api_errors.ErrChatNotFound
		}
		return err
	}
	return nil
}

func (s *ChatService) ChangeChat(caller dto.UserDTO, chatId int64, request dto.ChangeChatRequest) error {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return api_errors.ErrUnauthorized
	}

	chat, err := s.ChatRepository.GetById(chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return api_errors.ErrChatNotFound
		}
		return err
	}
	if chat.OwnerID != caller.ID {
		return api_errors.ErrNotEnoughPermissionsForChangeChat
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

	err = s.ChatRepository.UpdateById(chatId, updateData)
	if err != nil {
		if errors.As(err, &repositories.ErrDuplicate) {
			return api_errors.ErrChatAlreadyExists
		}
	}
	return nil
}

func (s *ChatService) GetListForUser(caller dto.UserDTO, page int) ([]dto.ChatDTO, error) {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return []dto.ChatDTO{}, nil
	}

	if page < 1 {
		return []dto.ChatDTO{}, api_errors.ErrInvalidPage
	}

	list, err := s.ChatRepository.GetListForUser(caller.ID, s.App.Config.Pagination.ChatList, (page-1)*s.App.Config.Pagination.ChatList)
	if err != nil {
		if errors.As(err, &repositories.ErrLimitMustBePositive) || errors.As(err, &repositories.ErrOffsetMustBePositive) {
			return []dto.ChatDTO{}, api_errors.ErrInvalidPage
		}
		return []dto.ChatDTO{}, err
	}

	result := make([]dto.ChatDTO, len(list))

	for i, v := range list {
		result[i] = v.ToDTO()
	}
	return result, nil
}

func (s *ChatService) Search(caller dto.UserDTO, name string, page int) ([]dto.ChatDTO, error) {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return []dto.ChatDTO{}, nil
	}

	if page < 1 {
		return []dto.ChatDTO{}, api_errors.ErrInvalidPage
	}

	list, err := s.ChatRepository.SearchForUser(caller.ID, name, s.App.Config.Pagination.ChatList, (page-1)*s.App.Config.Pagination.ChatList)
	if err != nil {
		if errors.As(err, &repositories.ErrLimitMustBePositive) || errors.As(err, &repositories.ErrOffsetMustBePositive) {
			return []dto.ChatDTO{}, api_errors.ErrInvalidPage
		}
		return []dto.ChatDTO{}, err
	}

	result := make([]dto.ChatDTO, len(list))

	for i, v := range list {
		result[i] = v.ToDTO()
	}
	return result, nil
}
