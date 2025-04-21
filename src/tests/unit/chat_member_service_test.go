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
	usecase_errors "libs/src/internal/usecase/errors"
	"reflect"
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
			expectedResp:  usecase_errors.AlreadyExistsError{},
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
				assert.Equal(t, reflect.TypeOf(tc.expectedResp), reflect.TypeOf(err))
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
			expectedResp:     usecase_errors.BadRequestError{},
			mustErr:          true,
		},
		{
			testName:  "Not enough permissions for inviting",
			inviterId: 1,
			GetMemberInfoResp: dto.MemberInfo{
				MemberRole: enums.MEMBER,
			},
			expectedResp: usecase_errors.PermissionError{},
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
			expectedResp:     usecase_errors.NotFoundError{},
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
			expectedResp: usecase_errors.NotFoundError{},
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
			expectedResp: usecase_errors.NotFoundError{},
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

			err := service.InviteToChat(&dto.UserDTO{ID: tc.inviterId, Role: enums.USER, IsActive: true}, tc.inviteeUsername, tc.chatId)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectedResp), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestChangeMemberRole(t *testing.T) {
	mockApp := GetAppMock()
	service := services.ChatMemberService{
		App: mockApp,
	}

	testCases := []struct {
		testName string

		targetName string
		callerId   int64
		chatId     int64
		memberId   int64
		newRole    string

		GetMemberInfoCallerResp dto.MemberInfo
		GetMemberInfoCallerErr  error

		GetByUsernameResp domain.User
		GetByUsernameErr  error

		GetMemberInfoTargetResp dto.MemberInfo
		GetMemberInfoTargetErr  error

		SetNewRoleResp error

		expectedResp error
		mustErr      bool
	}{
		{
			testName:     "An attempt to appoint the owner",
			callerId:     1,
			targetName:   "userTest",
			chatId:       1,
			memberId:     2,
			newRole:      "owner",
			expectedResp: usecase_errors.PermissionError{},
			mustErr:      true,
		},
		{
			testName:     "Role is not exists",
			callerId:     1,
			targetName:   "userTest",
			chatId:       1,
			memberId:     2,
			newRole:      "not_exists",
			expectedResp: usecase_errors.BadRequestError{},
			mustErr:      true,
		},
		{
			testName:               "Caller not found",
			callerId:               1,
			targetName:             "userTest",
			chatId:                 1,
			memberId:               2,
			newRole:                "admin",
			GetMemberInfoCallerErr: repositories.ErrRecordNotFound,
			expectedResp:           usecase_errors.NotFoundError{},
			mustErr:                true,
		},
		{
			testName:   "Caller has no permissions for changing role",
			callerId:   1,
			targetName: "userTest",
			chatId:     1,
			memberId:   2,
			newRole:    "admin",
			GetMemberInfoCallerResp: dto.MemberInfo{
				MemberRole: enums.MEMBER,
			},
			expectedResp: usecase_errors.PermissionError{},
			mustErr:      true,
		},
		{
			testName:   "Target not found",
			callerId:   1,
			targetName: "userTest",
			chatId:     1,
			memberId:   2,
			newRole:    "admin",
			GetMemberInfoCallerResp: dto.MemberInfo{
				MemberRole: enums.OWNER,
			},
			GetByUsernameErr: repositories.ErrRecordNotFound,
			expectedResp:     usecase_errors.NotFoundError{},
			mustErr:          true,
		},
		{
			testName:   "Target is anonymous",
			callerId:   1,
			targetName: "userTest",
			chatId:     1,
			memberId:   2,
			newRole:    "admin",
			GetMemberInfoCallerResp: dto.MemberInfo{
				MemberRole: enums.OWNER,
			},
			GetByUsernameResp: domain.User{
				Username: "userTest",
				IsActive: true,
				Role:     enums.ANONYMOUS,
			},
			expectedResp: usecase_errors.NotFoundError{},
			mustErr:      true,
		},
		{
			testName:   "Target is not active",
			callerId:   1,
			targetName: "userTest",
			chatId:     1,
			memberId:   2,
			newRole:    "admin",
			GetMemberInfoCallerResp: dto.MemberInfo{
				MemberRole: enums.OWNER,
			},
			GetByUsernameResp: domain.User{
				Username: "userTest",
				IsActive: false,
				Role:     enums.USER,
			},
			expectedResp: usecase_errors.NotFoundError{},
			mustErr:      true,
		},
		{
			testName:   "Target not in chat",
			callerId:   1,
			targetName: "userTest",
			chatId:     1,
			memberId:   2,
			newRole:    "admin",
			GetMemberInfoCallerResp: dto.MemberInfo{
				MemberRole: enums.OWNER,
			},
			GetByUsernameResp: domain.User{
				Username: "userTest",
				IsActive: true,
				Role:     enums.USER,
			},
			GetMemberInfoTargetErr: repositories.ErrRecordNotFound,
			expectedResp:           usecase_errors.BadRequestError{},
			mustErr:                true,
		},
		{
			testName:   "Target already has this role",
			callerId:   1,
			chatId:     1,
			targetName: "userTest",
			memberId:   2,
			newRole:    "admin",
			GetMemberInfoCallerResp: dto.MemberInfo{
				MemberRole: enums.OWNER,
			},
			GetByUsernameResp: domain.User{
				Username: "userTest",
				IsActive: true,
				Role:     enums.USER,
			},
			GetMemberInfoTargetResp: dto.MemberInfo{
				MemberRole: enums.ADMIN,
			},
			mustErr: false,
		},
		{
			testName:   "Success",
			callerId:   1,
			chatId:     1,
			memberId:   2,
			targetName: "userTest",
			newRole:    "admin",
			GetMemberInfoCallerResp: dto.MemberInfo{
				MemberRole: enums.OWNER,
			},
			GetByUsernameResp: domain.User{
				Username: "userTest",
				IsActive: true,
				Role:     enums.USER,
			},
			GetMemberInfoTargetResp: dto.MemberInfo{
				MemberRole: enums.MEMBER,
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
			mockChatMemberRepo.EXPECT().GetMemberInfo(tc.callerId, tc.chatId).Return(tc.GetMemberInfoCallerResp, tc.GetMemberInfoCallerErr)
			mockUserRepo.EXPECT().GetByUsername(mock.Anything).Return(tc.GetByUsernameResp, tc.GetByUsernameErr)
			mockChatMemberRepo.EXPECT().GetMemberInfo(mock.Anything, mock.Anything).Return(tc.GetMemberInfoTargetResp, tc.GetMemberInfoTargetErr)

			mockChatMemberRepo.EXPECT().SetNewRole(mock.Anything, mock.Anything, mock.Anything).Return(tc.SetNewRoleResp)

			err := service.ChangeMemberRole(dto.UserDTO{ID: tc.callerId}, tc.chatId, tc.targetName, tc.newRole)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectedResp), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestGetMemberList(t *testing.T) {
	mockApp := GetAppMock()
	service := services.ChatMemberService{
		App: mockApp,
	}

	testCases := []struct {
		testName          string
		caller            dto.UserDTO
		chatId            int64
		page              int
		searchName        string
		FilterResp        []domain.ChatMember
		FilterErr         error
		GetMemberListResp []dto.MemberPreview
		GetMemberListErr  error
		expectedResp      dto.MemberListPreview
		expectedErr       error
		mustErr           bool
	}{
		{
			testName: "User is not active",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: false,
			},
			expectedErr: usecase_errors.UnauthorizedError{},
			mustErr:     true,
		},
		{
			testName: "User is anonymous",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			expectedErr: usecase_errors.UnauthorizedError{},
			mustErr:     true,
		},
		{
			testName: "User is not in chat",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			chatId:      1,
			page:        1,
			expectedErr: usecase_errors.BadRequestError{},
			mustErr:     true,
		},
		{
			testName: "Db error",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			chatId: 1,
			page:   1,
			FilterResp: []domain.ChatMember{
				{
					ChatID:     1,
					UserID:     1,
					MemberRole: enums.MEMBER,
				},
			},
			GetMemberListErr: repositories.ErrOffsetMustBePositive,
			expectedErr:      usecase_errors.BadRequestError{},
			mustErr:          true,
		},
		{
			testName: "Success",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			chatId: 1,
			page:   1,
			FilterResp: []domain.ChatMember{
				{
					ChatID:     1,
					UserID:     1,
					MemberRole: enums.MEMBER,
				},
			},
			expectedResp: dto.MemberListPreview{
				Members: []dto.MemberPreview{},
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockChatMemberRepo := new(mocks.IChatMemberRepository)
		service.ChatMemberRepository = mockChatMemberRepo

		t.Run(tc.testName, func(t *testing.T) {
			mockChatMemberRepo.EXPECT().Filter(mock.Anything, mock.Anything, mock.Anything).Maybe().Return(tc.FilterResp, tc.FilterErr)
			mockChatMemberRepo.EXPECT().GetMembersPreview(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return(tc.GetMemberListResp, tc.GetMemberListErr)

			resp, err := service.GetList(tc.caller, tc.chatId, tc.page, tc.searchName)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectedErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedResp.Members), len(resp.Members))
			}
		})
	}
}
