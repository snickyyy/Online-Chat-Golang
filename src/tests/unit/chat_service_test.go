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
