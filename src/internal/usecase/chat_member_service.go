package services

import (
	"errors"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	usecase_errors "libs/src/internal/usecase/errors"
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
		return usecase_errors.AlreadyExistsError{Msg: "User already exists in chat"}
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
		return usecase_errors.UnauthorizedError{Msg: "You must be logged in to invite someone"}
	}
	inviterInfo, err := s.ChatMemberRepository.GetMemberInfo(inviter.ID, chatId)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return usecase_errors.BadRequestError{Msg: "Inviter is not a member of the chat"}
		}
		return err
	}
	if inviterInfo.MemberRole < enums.CHAT_ADMIN {
		return usecase_errors.PermissionError{Msg: "You have no permission to invite someone to the chat"}
	}

	invitee, err := s.UserRepository.GetByUsername(inviteeUsername)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.NotFoundError{Msg: "Invitee not found"}
		}
		return err
	}
	if invitee.Role == enums.ANONYMOUS || !invitee.IsActive {
		return usecase_errors.NotFoundError{Msg: "Invitee not found"}
	}

	err = s.CreateMember(invitee.ID, chatId)

	return err
}

func (s *ChatMemberService) ChangeMemberRole(caller dto.UserDTO, chatId int64, targetUsername string, newRole string) error {
	role, ex := enums.ChatLabelsToRoles[strings.ToLower(newRole)]
	if !ex {
		return usecase_errors.BadRequestError{Msg: "Invalid role"}
	}
	if role >= enums.OWNER {
		return usecase_errors.PermissionError{Msg: "You do not have permission to change role"}
	}

	callerInfo, err := s.ChatMemberRepository.GetMemberInfo(caller.ID, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.NotFoundError{Msg: "You are not a member of the chat"}
		}
		return err
	}

	if callerInfo.MemberRole < enums.OWNER {
		return usecase_errors.PermissionError{Msg: "You do not have permission to change role"}
	}

	target, err := s.UserRepository.GetByUsername(targetUsername)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.NotFoundError{Msg: "Target user not found"}
		}
		return err
	}

	if target.Role == enums.ANONYMOUS || !target.IsActive {
		return usecase_errors.NotFoundError{Msg: "Target user not found"}
	}

	targetInfo, err := s.ChatMemberRepository.GetMemberInfo(target.ID, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.BadRequestError{Msg: "Target user not in chat"}
		}
		return err
	}

	if targetInfo.MemberRole == byte(role) {
		return nil
	}

	err = s.ChatMemberRepository.SetNewRole(chatId, targetInfo.MemberID, byte(role))
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.BadRequestError{Msg: "Target user not in chat"}
		}
		return err
	}
	return nil
}

func (s *ChatMemberService) DeleteMember(caller dto.UserDTO, chatId int64, targetUsername string) error {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return usecase_errors.UnauthorizedError{Msg: "You must be logged in to delete someone"}
	}

	callerInfo, err := s.ChatMemberRepository.GetMemberInfo(caller.ID, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.BadRequestError{Msg: "You are not a member of the chat"}
		}
		return err
	}

	if callerInfo.MemberRole < enums.CHAT_ADMIN {
		return usecase_errors.PermissionError{Msg: "You do not have permission to delete someone"}
	}

	target, err := s.UserRepository.GetByUsername(targetUsername)
	if err != nil {
		if target.Role == enums.ANONYMOUS || !target.IsActive {
			return usecase_errors.NotFoundError{Msg: "Target user not found"}
		}
		return err
	}

	targetInfo, err := s.ChatMemberRepository.GetMemberInfo(target.ID, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.BadRequestError{Msg: "Target user not in chat"}
		}
		return err
	}

	if targetInfo.MemberRole >= callerInfo.MemberRole {
		return usecase_errors.PermissionError{Msg: "You do not have permission to delete target"}
	}

	err = s.ChatMemberRepository.DeleteMember(targetInfo.MemberID, chatId)
	if err != nil {
		if errors.As(err, &repositories.ErrRecordNotFound) {
			return usecase_errors.BadRequestError{Msg: "Target user not in chat"}
		}
		return err
	}
	return nil
}

func (s *ChatMemberService) GetList(caller dto.UserDTO, chatId int64, page int, searchName string) (dto.MemberListPreview, error) {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return dto.MemberListPreview{}, usecase_errors.UnauthorizedError{Msg: "You must be logged in to get members"}
	}

	if page < 1 {
		return dto.MemberListPreview{}, usecase_errors.BadRequestError{Msg: "Invalid page"}
	}

	callerMember, err := s.ChatMemberRepository.Filter("chat_id = ? AND user_id = ?", chatId, caller.ID)
	if err != nil {
		return dto.MemberListPreview{}, err
	}

	if len(callerMember) != 1 {
		return dto.MemberListPreview{}, usecase_errors.BadRequestError{Msg: "You are not a member of the chat"}
	}

	res, err := s.ChatMemberRepository.GetMembersPreview(chatId, 25, (page-1)*25, searchName)
	if err != nil {
		if errors.As(err, &repositories.ErrLimitMustBePositive) || errors.As(err, &repositories.ErrOffsetMustBePositive) {
			return dto.MemberListPreview{}, usecase_errors.BadRequestError{Msg: "Invalid page"}
		}
		return dto.MemberListPreview{}, err
	}

	return dto.MemberListPreview{Members: res}, nil
}
