package unit

import (
	"context"
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

func TestCreateChat(t *testing.T) {
	mockApp := GetAppMock()

	chatService := services.ChatService{
		App: mockApp,
	}

	testCases := []struct {
		testName string
		request  dto.CreateChatRequest
		user     dto.UserDTO
		mockResp error
		mustErr  bool
		resp     dto.ChatDTO
		respErr  error
	}{
		{
			"TestCreateChatUnauthorized",
			dto.CreateChatRequest{
				Title:       "Test Chat",
				Description: "Test Description",
			},
			dto.UserDTO{
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			usecase_errors.UnauthorizedError{},
			true,
			dto.ChatDTO{},
			usecase_errors.UnauthorizedError{},
		},
		{
			"TestCreateChatDuplicate",
			dto.CreateChatRequest{
				Title:       "Test Chat",
				Description: "Test Description",
			},
			dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			repositories.ErrDuplicate,
			true,
			dto.ChatDTO{},
			usecase_errors.AlreadyExistsError{},
		},
		{
			"TestCreateChatSuccess",
			dto.CreateChatRequest{
				Title:       "Test Chat",
				Description: "Test Description",
			},
			dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			nil,
			false,
			dto.ChatDTO{
				Title:       "Test Chat",
				Description: "Test Description",
			},
			nil,
		},
	}
	for _, tc := range testCases {
		MockChatRepository := new(mocks.IChatRepository)
		chatService.ChatRepository = MockChatRepository

		t.Run(tc.testName, func(t *testing.T) {
			MockChatRepository.EXPECT().Create(context.Background(), mock.Anything).Return(tc.mockResp)

			chat, err := chatService.CreateChat(context.Background(), tc.request, tc.user)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(tc.respErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, chat.Title, tc.request.Title)
				assert.Equal(t, chat.Description, tc.request.Description)
			}
		})
	}
}
func TestDeleteChat(t *testing.T) {
	mockApp := GetAppMock()
	service := services.ChatService{
		App: mockApp,
	}

	testCases := []struct {
		testName    string
		caller      dto.UserDTO
		chatID      int64
		GetByIdResp domain.Chat
		GetByIdErr  error
		DeleteResp  error
		expectErr   error
		mustErr     bool
	}{
		{
			testName: "TestDeleteChatNotFound",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			chatID:      1,
			GetByIdResp: domain.Chat{},
			GetByIdErr:  repositories.ErrRecordNotFound,
			expectErr:   usecase_errors.NotFoundError{},
			mustErr:     true,
		},
		{
			testName: "TestDeleteChatCallerIsNotOwner",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			chatID: 1,
			GetByIdResp: domain.Chat{
				OwnerID: 12,
			},
			expectErr: usecase_errors.PermissionError{},
			mustErr:   true,
		},
		{
			testName: "TestDeleteChatSuccess",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			chatID: 1,
			GetByIdResp: domain.Chat{
				OwnerID: 1,
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockChatRepository := new(mocks.IChatRepository)
		service.ChatRepository = mockChatRepository

		t.Run(tc.testName, func(t *testing.T) {
			mockChatRepository.EXPECT().GetById(context.Background(), tc.chatID).Return(tc.GetByIdResp, tc.GetByIdErr)
			mockChatRepository.EXPECT().DeleteById(context.Background(), tc.chatID).Maybe().Return(tc.DeleteResp)

			err := service.DeleteChat(context.Background(), tc.caller, tc.chatID)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
			}
		})

	}
}
func TestGetChatListForUser(t *testing.T) {
	mockApp := GetAppMock()
	chatService := services.ChatService{
		App: mockApp,
	}

	testCases := []struct {
		testName   string
		caller     dto.UserDTO
		page       int
		RepoResp   []domain.Chat
		RepoErr    error
		expectResp []dto.ChatDTO
		expectErr  error
		mustErr    bool
	}{
		{
			testName: "TestGetChatListForUserUnauthorized",
			caller: dto.UserDTO{
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			page:    1,
			mustErr: false,
		},
		{
			testName: "TestGetChatListForUserNotActive",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: false,
			},
			page:    1,
			mustErr: false,
		},
		{
			testName: "TestGetChatListForUserInvalidPage",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			page:      -1,
			expectErr: usecase_errors.BadRequestError{},
			mustErr:   true,
		},
		{
			testName: "TestGetChatListForUserDbError",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			page:      1,
			RepoErr:   repositories.ErrLimitMustBePositive,
			expectErr: usecase_errors.BadRequestError{},
			mustErr:   true,
		},
		{
			testName: "TestGetChatListForUserSuccess",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			page: 1,
			RepoResp: []domain.Chat{
				{
					Title:       "Test Chat 1",
					Description: "Test Description 1",
					OwnerID:     1,
				}, {
					Title:       "Test Chat 2",
					Description: "Test Description 1",
					OwnerID:     2,
				},
				{
					Title:       "Test Chat 3",
					Description: "Test Description 1",
					OwnerID:     123,
				}, {
					Title:       "Test Chat 4",
					Description: "Test Description 1",
					OwnerID:     18,
				},
			},
			expectResp: []dto.ChatDTO{
				{
					Title:       "Test Chat 1",
					Description: "Test Description 1",
					OwnerID:     1,
				},
				{
					Title:       "Test Chat 2",
					Description: "Test Description 1",
					OwnerID:     2,
				},
				{
					Title:       "Test Chat 3",
					Description: "Test Description 1",
					OwnerID:     123,
				},
				{
					Title:       "Test Chat 4",
					Description: "Test Description 1",
					OwnerID:     18,
				},
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockChatRepository := new(mocks.IChatRepository)
		chatService.ChatRepository = mockChatRepository

		t.Run(tc.testName, func(t *testing.T) {
			mockChatRepository.EXPECT().GetListForUser(context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(tc.RepoResp, tc.RepoErr)

			chats, err := chatService.GetListForUser(context.Background(), tc.caller, tc.page)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(chats), len(tc.expectResp))
			}
		})
	}
}
func TestSearchChat(t *testing.T) {
	mockApp := GetAppMock()
	service := services.ChatService{
		App: mockApp,
	}

	testCases := []struct {
		testName   string
		caller     dto.UserDTO
		query      string
		page       int
		RepoResp   []domain.Chat
		RepoErr    error
		expectResp []dto.ChatDTO
		expectErr  error
		mustErr    bool
	}{
		{
			testName: "TestSearchChatUnauthorized",
			caller: dto.UserDTO{
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			query:   "Test",
			page:    1,
			mustErr: false,
		},
		{
			testName: "TestSearchChatNotActive",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: false,
			},
			query:   "Test",
			page:    1,
			mustErr: false,
		},
		{
			testName: "TestSearchChatInvalidPage",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			query:     "Test",
			page:      -1,
			expectErr: usecase_errors.BadRequestError{},
			mustErr:   true,
		},
		{
			testName: "TestSearchChatDbError",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			query:     "Test",
			page:      1,
			RepoErr:   repositories.ErrLimitMustBePositive,
			expectErr: usecase_errors.BadRequestError{},
			mustErr:   true,
		},
		{
			testName: "TestSearchChatSuccess",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			query: "Test",
			page:  1,
			RepoResp: []domain.Chat{
				{
					Title:       "Test Chat 1",
					Description: "Test Description 1",
					OwnerID:     1,
				}, {
					Title:       "Test Chat 2",
					Description: "Test Description 1",
					OwnerID:     2,
				},
				{
					Title:       "Test Chat 3",
					Description: "Test Description 1",
					OwnerID:     123,
				}, {
					Title:       "Test Chat 4",
					Description: "Test Description 1",
					OwnerID:     18,
				},
			},
			expectResp: []dto.ChatDTO{
				{
					Title:       "Test Chat 1",
					Description: "Test Description 1",
					OwnerID:     1,
				},
				{
					Title:       "Test Chat 2",
					Description: "Test Description 1",
					OwnerID:     2,
				},
				{
					Title:       "Test Chat 3",
					Description: "Test Description 1",
					OwnerID:     123,
				},
				{
					Title:       "Test Chat 4",
					Description: "Test Description 1",
					OwnerID:     18,
				},
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockChatRepository := new(mocks.IChatRepository)
		service.ChatRepository = mockChatRepository

		t.Run(tc.testName, func(t *testing.T) {
			mockChatRepository.EXPECT().SearchForUser(context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return(tc.RepoResp, tc.RepoErr)

			chats, err := service.Search(context.Background(), tc.caller, tc.query, tc.page)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(chats), len(tc.expectResp))
			}
		})
	}
}
func TestGetChatById(t *testing.T) {
	mockApp := GetAppMock()
	service := services.ChatService{
		App: mockApp,
	}

	testCases := []struct {
		testName    string
		caller      dto.UserDTO
		chatId      int64
		FilterResp  []domain.ChatMember
		FilterErr   error
		GetByIdResp domain.Chat
		GetByIdErr  error
		expectResp  dto.ChatDTO
		expectErr   error
		mustErr     bool
	}{
		{
			testName: "TestGetChatByIdUnauthorized",
			caller: dto.UserDTO{
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			chatId:    1,
			expectErr: usecase_errors.UnauthorizedError{},
			mustErr:   true,
		},
		{
			testName: "TestGetChatByIdNotActive",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: false,
			},
			chatId:    1,
			expectErr: usecase_errors.UnauthorizedError{},
			mustErr:   true,
		},
		{
			testName: "TestGetChatByIdUserNotInChat",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			chatId:    1,
			expectErr: usecase_errors.NotFoundError{},
			mustErr:   true,
		},
		{
			testName: "TestGetChatByIdChatNotExist",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			chatId: 1,
			FilterResp: []domain.ChatMember{
				{
					ChatID: 1,
					UserID: 1,
				},
			},
			GetByIdErr: repositories.ErrRecordNotFound,
			expectErr:  usecase_errors.NotFoundError{},
			mustErr:    true,
		},
		{
			testName: "TestGetChatByIdSuccess",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			chatId: 1,
			FilterResp: []domain.ChatMember{
				{
					ChatID: 1,
					UserID: 1,
				},
			},
			GetByIdResp: domain.Chat{
				Title:       "Test Chat 1",
				Description: "Test Description 1",
				OwnerID:     1,
			},
			expectResp: dto.ChatDTO{
				Title:       "Test Chat 1",
				Description: "Test Description 1",
				OwnerID:     1,
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockChatRepository := new(mocks.IChatRepository)
		mockChatMemberRepository := new(mocks.IChatMemberRepository)
		service.ChatRepository = mockChatRepository
		service.ChatMemberRepository = mockChatMemberRepository

		t.Run(tc.testName, func(t *testing.T) {
			mockChatMemberRepository.EXPECT().Filter(context.Background(), mock.Anything, mock.Anything, mock.Anything).Maybe().Return(tc.FilterResp, tc.FilterErr)
			mockChatRepository.EXPECT().GetById(context.Background(), mock.Anything).Maybe().Return(tc.GetByIdResp, tc.GetByIdErr)

			chat, err := service.GetById(context.Background(), tc.caller, tc.chatId)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, chat, tc.expectResp)
			}
		})
	}
}
