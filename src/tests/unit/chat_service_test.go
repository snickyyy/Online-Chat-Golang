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

func TestCreateChatUnauthorized(t *testing.T) {
	mockApp := GetAppMock()
	MockChatRepository := new(mocks.IChatRepository)

	MockChatRepository.EXPECT().Create(mock.Anything).Return(nil)

	chatService := services.ChatService{
		App:            mockApp,
		ChatRepository: MockChatRepository,
	}

	request := dto.CreateChatRequest{
		Title:       "Test Chat",
		Description: "Test Description",
	}
	testCases := []struct {
		user  dto.UserDTO
		error error
	}{
		{
			user: dto.UserDTO{
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			error: api_errors.ErrUnauthorized,
		},
		{
			user: dto.UserDTO{
				Role:     enums.ADMIN,
				IsActive: false,
			},
			error: api_errors.ErrUnauthorized,
		},
		{
			user: dto.UserDTO{
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			error: api_errors.ErrUnauthorized,
		},
	}
	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]
		_, err := chatService.CreateChat(request, tc.user)
		assert.Equal(t, err, tc.error, "I: %d", i)
	}
}
func TestCreateChatDuplicate(t *testing.T) {
	mockApp := GetAppMock()
	MockChatRepository := new(mocks.IChatRepository)

	MockChatRepository.EXPECT().Create(mock.Anything).Return(repositories.ErrDuplicate)
	chatService := services.ChatService{App: mockApp, ChatRepository: MockChatRepository}

	request := dto.CreateChatRequest{
		Title:       "Test Chat",
		Description: "Test Description",
	}
	user := dto.UserDTO{
		ID:       1,
		Role:     enums.USER,
		IsActive: true,
	}
	_, err := chatService.CreateChat(request, user)
	assert.Equal(t, err, api_errors.ErrChatAlreadyExists)
}
func TestCreateChatSuccess(t *testing.T) {
	mockApp := GetAppMock()
	MockChatRepository := new(mocks.IChatRepository)

	MockChatRepository.EXPECT().Create(mock.Anything).Return(nil)
	chatService := services.ChatService{App: mockApp, ChatRepository: MockChatRepository}

	request := dto.CreateChatRequest{
		Title:       "Test Chat",
		Description: "Test Description",
	}
	user := dto.UserDTO{
		ID:       1,
		Role:     enums.USER,
		IsActive: true,
	}
	chat, err := chatService.CreateChat(request, user)
	assert.NoError(t, err)
	assert.Equal(t, chat.Title, request.Title)
	assert.Equal(t, chat.Description, request.Description)
	assert.Equal(t, chat.OwnerID, user.ID)
}
