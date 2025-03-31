package services

import (
	"errors"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
	"strings"
)

type ChatService struct {
	App                  *settings.App
	UserRepository       *repositories.UserRepository
	ChatRepository       *repositories.ChatRepository
	ChatMemberRepository *repositories.ChatMemberRepository
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

func (s *ChatService) InviteToChat(inviter *dto.UserDTO, inviteeUsername string, chatId int64) error {
	inviterInfo, err := s.ChatRepository.GetMemberInfo(inviter.ID, chatId)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return api_errors.ErrInviterNotInChat
		}
		return err
	}
	if inviterInfo.MemberRole < enums.CHAT_ADMIN {
		return api_errors.ErrNotEnoughPermissionsForInviting
	}

	invitee, err := s.UserRepository.GetByUsername(inviteeUsername)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return api_errors.ErrUserNotFound
		}
		return err
	}
	if invitee.Role == enums.ANONYMOUS || !invitee.IsActive {
		return api_errors.ErrUserNotFound
	}

	memberCount, err := s.ChatMemberRepository.Count("chat_id = ? AND user_id = ?", chatId, invitee.ID)
	if err != nil {
		return err
	}
	if memberCount > 0 {
		return api_errors.ErrUserAlreadyInChat
	}

	newMember := domain.ChatMember{
		ChatID:     chatId,
		UserID:     invitee.ID,
		MemberRole: enums.MEMBER,
	}
	_, err = s.ChatMemberRepository.Create(&newMember)

	return err
}

func (s *ChatService) ChangeMemberRole(caller dto.UserDTO, chatId int64, targetUsername string, newRole string) error {
	role, ex := enums.ChatLabelsToRoles[strings.ToLower(newRole)]
	if !ex {
		return api_errors.ErrInvalidData
	}
	if !(role < enums.OWNER) {
		return api_errors.ErrNotEnoughPermissionsForChangeRole
	}

	callerInfo, err := s.ChatRepository.GetMemberInfo(caller.ID, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return api_errors.ErrUserNotInChat
		}
		return err
	}

	if callerInfo.MemberRole != enums.OWNER {
		return api_errors.ErrNotEnoughPermissionsForChangeRole
	}

	target, err := s.UserRepository.GetByUsername(targetUsername)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return api_errors.ErrUserNotFound
		}
		return err
	}

	if target.Role == enums.ANONYMOUS || !target.IsActive {
		return api_errors.ErrUserNotFound
	}

	targetInfo, err := s.ChatRepository.GetMemberInfo(target.ID, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return api_errors.ErrUserNotInChat
		}
		return err
	}

	if targetInfo.MemberRole == byte(role) {
		return nil
	}

	err = s.ChatMemberRepository.SetNewRole(chatId, targetInfo.MemberID, byte(role))
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return api_errors.ErrUserNotFound
		}
		return err
	}
	return nil
}
