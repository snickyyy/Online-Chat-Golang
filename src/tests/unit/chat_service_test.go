package unit

import (
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
			api_errors.ErrUnauthorized,
			true,
			dto.ChatDTO{},
			api_errors.ErrUnauthorized,
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
			api_errors.ErrChatAlreadyExists,
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
			MockChatRepository.EXPECT().Create(mock.Anything).Return(tc.mockResp)

			chat, err := chatService.CreateChat(tc.request, tc.user)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, err, tc.respErr)
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
			expectErr:   api_errors.ErrChatNotFound,
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
			expectErr: api_errors.ErrNotEnoughPermissionsForDelete,
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
			mockChatRepository.EXPECT().GetById(tc.chatID).Return(tc.GetByIdResp, tc.GetByIdErr)
			mockChatRepository.EXPECT().DeleteById(tc.chatID).Maybe().Return(tc.DeleteResp)

			err := service.DeleteChat(tc.caller, tc.chatID)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectErr, err)
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
			expectErr: api_errors.ErrInvalidPage,
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
			expectErr: api_errors.ErrInvalidPage,
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
			mockChatRepository.EXPECT().GetListForUser(mock.Anything, mock.Anything, mock.Anything).Return(tc.RepoResp, tc.RepoErr)

			chats, err := chatService.GetListForUser(tc.caller, tc.page)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectErr, err)
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
			expectErr: api_errors.ErrInvalidPage,
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
			expectErr: api_errors.ErrInvalidPage,
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
			mockChatRepository.EXPECT().SearchForUser(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return(tc.RepoResp, tc.RepoErr)

			chats, err := service.Search(tc.caller, tc.query, tc.page)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectErr, err)
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
			expectErr: api_errors.ErrUnauthorized,
			mustErr:   true,
		},
		{
			testName: "TestGetChatByIdNotActive",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: false,
			},
			chatId:    1,
			expectErr: api_errors.ErrUnauthorized,
			mustErr:   true,
		},
		{
			testName: "TestGetChatByIdUserNotInChat",
			caller: dto.UserDTO{
				Role:     enums.USER,
				IsActive: true,
			},
			chatId:    1,
			expectErr: api_errors.ErrChatNotFound,
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
			expectErr:  api_errors.ErrChatNotFound,
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
			mockChatMemberRepository.EXPECT().Filter(mock.Anything, mock.Anything, mock.Anything).Maybe().Return(tc.FilterResp, tc.FilterErr)
			mockChatRepository.EXPECT().GetById(mock.Anything).Maybe().Return(tc.GetByIdResp, tc.GetByIdErr)

			chat, err := service.GetById(tc.caller, tc.chatId)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, chat, tc.expectResp)
			}
		})
	}
}
