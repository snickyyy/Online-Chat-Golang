package unit

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"libs/src/internal/domain/enums"
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
