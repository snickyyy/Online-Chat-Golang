package unit

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/mocks"
	"libs/src/internal/repositories"
	services "libs/src/internal/usecase"
	api_errors "libs/src/internal/usecase/errors"
	"testing"
)

func TestCreateMember(t *testing.T) {
	mockApp := GetAppMock()
	service := services.ChatMemberService{
		App: mockApp,
	}

	dbErr := errors.New("internal db err")

	testCases := []struct {
		testName       string
		userId         int64
		chatId         int64
		RepoCountResp  int64
		RepoCountErr   error
		RepoCreateResp error
		expectedResp   error
		mustErr        bool
	}{
		{
			testName:      "User already in chat",
			userId:        1,
			chatId:        1,
			RepoCountResp: 1,
			expectedResp:  api_errors.ErrUserAlreadyInChat,
			mustErr:       true,
		},
		{
			testName:     "DataBase error",
			userId:       1,
			chatId:       1,
			RepoCountErr: dbErr,
			expectedResp: dbErr,
			mustErr:      true,
		},
		{
			testName: "Success",
			userId:   1,
			chatId:   1,
			mustErr:  false,
		},
	}

	for _, tc := range testCases {
		mockRepository := new(mocks.IChatMemberRepository)
		service.ChatMemberRepository = mockRepository

		t.Run(tc.testName, func(t *testing.T) {
			mockRepository.EXPECT().Count(mock.Anything, mock.Anything, mock.Anything).Return(tc.RepoCountResp, tc.RepoCountErr)
			mockRepository.EXPECT().Create(mock.Anything).Return(tc.RepoCreateResp)

			err := service.CreateMember(tc.userId, tc.chatId)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedResp, err)
			} else {
				assert.NoError(t, err)

			}
		})
	}
}
func TestInviteToChat(t *testing.T) {
	mockApp := GetAppMock()
	service := services.ChatMemberService{
		App: mockApp,
	}

	testCases := []struct {
		testName string

		inviterId       int64
		inviteeUsername string
		chatId          int64

		GetMemberInfoResp dto.MemberInfo
		GetMemberInfoErr  error

		GetByUsernameResp domain.User
		GetByUsernameErr  error

		expectedResp error
		mustErr      bool
	}{
		{
			testName:         "Inviter not in chat",
			inviterId:        1,
			GetMemberInfoErr: repositories.ErrRecordNotFound,
			expectedResp:     api_errors.ErrInviterNotInChat,
			mustErr:          true,
		},
		{
			testName:  "Not enough permissions for inviting",
			inviterId: 1,
			GetMemberInfoResp: dto.MemberInfo{
				MemberRole: enums.MEMBER,
			},
			expectedResp: api_errors.ErrNotEnoughPermissionsForInviting,
			mustErr:      true,
		},
		{
			testName:        "Invitee not found",
			inviterId:       1,
			inviteeUsername: "invitee",
			chatId:          1,
			GetMemberInfoResp: dto.MemberInfo{
				MemberRole: enums.CHAT_ADMIN,
			},
			GetByUsernameErr: repositories.ErrRecordNotFound,
			expectedResp:     api_errors.ErrUserNotFound,
			mustErr:          true,
		},
		{
			testName:        "Invitee is anonymous",
			inviterId:       1,
			inviteeUsername: "invitee",
			chatId:          1,
			GetMemberInfoResp: dto.MemberInfo{
				MemberRole: enums.CHAT_ADMIN,
			},
			GetByUsernameResp: domain.User{
				Username: "invitee",
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			expectedResp: api_errors.ErrUserNotFound,
			mustErr:      true,
		},
		{
			testName:        "Invitee is not active",
			inviterId:       1,
			inviteeUsername: "invitee",
			chatId:          1,
			GetMemberInfoResp: dto.MemberInfo{
				MemberRole: enums.CHAT_ADMIN,
			},
			GetByUsernameResp: domain.User{
				Username: "invitee",
				Role:     enums.USER,
				IsActive: false,
			},
			expectedResp: api_errors.ErrUserNotFound,
			mustErr:      true,
		},
		{
			testName:        "Success",
			inviterId:       1,
			inviteeUsername: "invitee",
			chatId:          1,
			GetMemberInfoResp: dto.MemberInfo{
				MemberRole: enums.CHAT_ADMIN,
			},
			GetByUsernameResp: domain.User{
				Username: "invitee",
				Role:     enums.USER,
				IsActive: true,
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockChatMemberRepo := new(mocks.IChatMemberRepository)
		mockUserRepo := new(mocks.IUserRepository)
		service.ChatMemberRepository = mockChatMemberRepo
		service.UserRepository = mockUserRepo

		t.Run(tc.testName, func(t *testing.T) {
			mockChatMemberRepo.EXPECT().GetMemberInfo(mock.Anything, mock.Anything).Return(tc.GetMemberInfoResp, tc.GetMemberInfoErr)
			mockUserRepo.EXPECT().GetByUsername(mock.Anything).Return(tc.GetByUsernameResp, tc.GetByUsernameErr)

			// Mocking the CreateMember method
			mockChatMemberRepo.EXPECT().Count(mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
			mockChatMemberRepo.EXPECT().Create(mock.Anything).Return(nil)

			err := service.InviteToChat(&dto.UserDTO{ID: tc.inviterId}, tc.inviteeUsername, tc.chatId)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedResp, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
