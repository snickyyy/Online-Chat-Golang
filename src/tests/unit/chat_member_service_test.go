package unit

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"libs/src/internal/mocks"
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
