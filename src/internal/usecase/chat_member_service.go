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

type ChatMemberService struct {
	App                  *settings.App
	ChatMemberRepository repositories.IChatMemberRepository
	UserRepository       repositories.IUserRepository
	ChatRepository       repositories.IChatRepository
}

func NewChatMemberService(app *settings.App) *ChatMemberService {
	return &ChatMemberService{
		App:                  app,
		ChatMemberRepository: repositories.NewChatMemberRepository(app),
		UserRepository:       repositories.NewUserRepository(app),
		ChatRepository:       repositories.NewChatRepository(app),
	}
}

func (s *ChatMemberService) CreateMember(userId int64, chatId int64) error {
	memberCount, err := s.ChatMemberRepository.Count("chat_id = ? AND user_id = ?", chatId, userId)
	if err != nil {
		return err
	}
	if memberCount > 0 {
		return api_errors.ErrUserAlreadyInChat
	}

	member := domain.ChatMember{
		UserID:     userId,
		ChatID:     chatId,
		MemberRole: enums.MEMBER,
	}
	err = s.ChatMemberRepository.Create(&member)
	return err
}

func (s *ChatMemberService) InviteToChat(inviter *dto.UserDTO, inviteeUsername string, chatId int64) error {
	if inviter.Role == enums.ANONYMOUS || !inviter.IsActive {
		return api_errors.ErrUnauthorized
	}
	inviterInfo, err := s.ChatMemberRepository.GetMemberInfo(inviter.ID, chatId)
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

	err = s.CreateMember(invitee.ID, chatId)

	return err
}

func (s *ChatMemberService) ChangeMemberRole(caller dto.UserDTO, chatId int64, targetUsername string, newRole string) error {
	role, ex := enums.ChatLabelsToRoles[strings.ToLower(newRole)]
	if !ex {
		return api_errors.ErrInvalidData
	}
	if role >= enums.OWNER {
		return api_errors.ErrNotEnoughPermissionsForChangeRole
	}

	callerInfo, err := s.ChatMemberRepository.GetMemberInfo(caller.ID, chatId)
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

	targetInfo, err := s.ChatMemberRepository.GetMemberInfo(target.ID, chatId)
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
