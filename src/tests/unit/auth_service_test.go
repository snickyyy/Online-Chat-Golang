package unit

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"libs/src/internal/domain/enums"
	"libs/src/internal/dto"
	"libs/src/internal/mocks"
	services "libs/src/internal/usecase"
	api_errors "libs/src/internal/usecase/errors"
	"testing"
)

func TestCheckEmailToken(t *testing.T) {
	mockApp := GetAppMock()
	service := services.AuthService{
		App: mockApp,
	}

	testCases := []struct {
		testName  string
		sessionId string
		mockResp  dto.UserDTO
		mockErr   error
		respError error
		Resp      dto.UserDTO
		mustErr   bool
	}{
		{
			testName:  "CheckEmailTokenInvalidSession",
			sessionId: "invalidSession",
			mockResp:  dto.UserDTO{},
			mockErr:   api_errors.ErrInvalidToken,
			respError: api_errors.ErrInvalidToken,
			Resp:      dto.UserDTO{},
			mustErr:   true,
		},
		{
			testName:  "CheckEmailUserAlreadyLoggedIn",
			sessionId: "test",
			mockResp: dto.UserDTO{
				ID:       1,
				Username: "test",
				Email:    "test@ocg.com",
				IsActive: true,
				Role:     enums.USER,
			},
			mockErr:   nil,
			respError: api_errors.ErrInvalidToken,
			Resp:      dto.UserDTO{},
			mustErr:   true,
		},
		{
			testName:  "CheckEmailSuccess",
			sessionId: "test",
			mockResp: dto.UserDTO{
				ID:       1,
				Username: "test",
				Email:    "test@ocg.com",
				IsActive: false,
				Role:     enums.ANONYMOUS,
			},
			mockErr:   nil,
			respError: nil,
			Resp: dto.UserDTO{
				ID:       1,
				Username: "test",
				Email:    "test@ocg.com",
				IsActive: false,
				Role:     enums.ANONYMOUS,
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockSessionService := new(mocks.ISessionService)
		service.SessionService = mockSessionService

		t.Run(tc.testName, func(t *testing.T) {
			mockSessionService.EXPECT().GetUserByEmailSession(tc.sessionId).Return(tc.mockResp, tc.mockErr)

			res, err := service.CheckEmailToken(tc.sessionId)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.respError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Resp, res)
			}
		})
	}
}
func TestConfirmAccount(t *testing.T) {
	mockApp := GetAppMock()
	service := services.AuthService{
		App: mockApp,
	}

	testCases := []struct {
		testName                       string
		Param                          string
		SessServGetUserByEmailSessResp dto.UserDTO
		SessServGetUserByEmailSessErr  error
		UserRepoUpdateByIdResp         error
		SessServSetSessionResp         string
		SessServSetSessionErr          error
		ExpectedResp                   string
		ExpectedError                  error
		mustErr                        bool
	}{
		{
			testName:                       "ConfirmAccountInvalidSession",
			Param:                          "invalidSession",
			SessServGetUserByEmailSessResp: dto.UserDTO{},
			SessServGetUserByEmailSessErr:  api_errors.ErrInvalidSession,
			UserRepoUpdateByIdResp:         nil,
			SessServSetSessionResp:         "",
			SessServSetSessionErr:          nil,
			ExpectedResp:                   "",
			ExpectedError:                  api_errors.ErrInvalidToken,
			mustErr:                        true,
		},
		{
			testName: "ConfirmAccountUserAlreadyLoggedIn",
			Param:    "test",
			SessServGetUserByEmailSessResp: dto.UserDTO{
				ID:       1,
				Role:     enums.USER,
				IsActive: true,
			},
			SessServGetUserByEmailSessErr: nil,
			UserRepoUpdateByIdResp:        nil,
			SessServSetSessionResp:        "",
			SessServSetSessionErr:         nil,
			ExpectedResp:                  "",
			ExpectedError:                 api_errors.ErrInvalidToken,
			mustErr:                       true,
		},
		{
			testName: "ConfirmAccountSuccess",
			Param:    "test",
			SessServGetUserByEmailSessResp: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			SessServGetUserByEmailSessErr: nil,
			UserRepoUpdateByIdResp:        nil,
			SessServSetSessionResp:        "test",
			SessServSetSessionErr:         nil,
			ExpectedResp:                  "test",
			ExpectedError:                 nil,
			mustErr:                       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository := new(mocks.IUserRepository)
			mockSessionService := new(mocks.ISessionService)
			service.UserRepository = mockUserRepository
			service.SessionService = mockSessionService
			mockSessionService.EXPECT().GetUserByEmailSession(mock.Anything).Return(tc.SessServGetUserByEmailSessResp, tc.SessServGetUserByEmailSessErr)
			mockUserRepository.EXPECT().UpdateById(mock.Anything, mock.Anything).Maybe().Return(tc.UserRepoUpdateByIdResp)
			mockSessionService.EXPECT().SetSession(mock.Anything).Maybe().Return(tc.SessServSetSessionResp, tc.SessServSetSessionErr)
			mockSessionService.EXPECT().DeleteSession(mock.Anything, mock.Anything).Maybe().Return(tc.UserRepoUpdateByIdResp)

			res, err := service.ConfirmAccount(tc.Param)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.ExpectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedResp, res)
			}
		})
	}
}
