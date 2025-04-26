package unit

import (
	"context"
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
	"libs/src/pkg/utils"
	"reflect"
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
			mockErr:   usecase_errors.BadRequestError{Msg: "Invalid token"},
			respError: usecase_errors.BadRequestError{},
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
			respError: usecase_errors.BadRequestError{},
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
			mockSessionService.EXPECT().GetUserByEmailSession(context.Background(), tc.sessionId).Return(tc.mockResp, tc.mockErr)

			res, err := service.CheckEmailToken(context.Background(), tc.sessionId)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.respError), reflect.TypeOf(err))
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
		caller                         dto.UserDTO
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
			testName: "ConfirmAccountAlreadyActive",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			ExpectedError: usecase_errors.BadRequestError{},
			mustErr:       true,
		},
		{
			testName: "ConfirmAccountInvalidSession",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			Param:                          "invalidSession",
			SessServGetUserByEmailSessResp: dto.UserDTO{},
			SessServGetUserByEmailSessErr:  usecase_errors.BadRequestError{},
			UserRepoUpdateByIdResp:         nil,
			SessServSetSessionResp:         "",
			SessServSetSessionErr:          nil,
			ExpectedResp:                   "",
			ExpectedError:                  usecase_errors.BadRequestError{},
			mustErr:                        true,
		},
		{
			testName: "ConfirmAccountUserAlreadyLoggedIn",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			Param: "test",
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
			ExpectedError:                 usecase_errors.BadRequestError{},
			mustErr:                       true,
		},
		{
			testName: "ConfirmAccountSuccess",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			Param: "test",
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
			mockSessionService.EXPECT().GetUserByEmailSession(context.Background(), mock.Anything).Return(tc.SessServGetUserByEmailSessResp, tc.SessServGetUserByEmailSessErr)
			mockUserRepository.EXPECT().UpdateById(context.Background(), mock.Anything, mock.Anything).Maybe().Return(tc.UserRepoUpdateByIdResp)
			mockSessionService.EXPECT().SetSession(context.Background(), mock.Anything).Maybe().Return(tc.SessServSetSessionResp, tc.SessServSetSessionErr)
			mockSessionService.EXPECT().DeleteSession(context.Background(), mock.Anything, mock.Anything).Maybe().Return(tc.UserRepoUpdateByIdResp)

			res, err := service.ConfirmAccount(context.Background(), tc.caller, tc.Param)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.ExpectedError), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedResp, res)
			}
		})
	}
}
func TestRegisterUser(t *testing.T) {
	mockApp := GetAppMock()
	service := services.AuthService{
		App: mockApp,
	}

	testCases := []struct {
		testName        string
		caller          dto.UserDTO
		data            dto.RegisterRequest
		UserRepoResp    error
		SessionServResp string
		SessionServErr  error
		respError       error
		mustErr         bool
	}{
		{
			testName: "RegisterUserAlreadyLoggedIn",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			respError: usecase_errors.BadRequestError{},
			mustErr:   true,
		},
		{
			testName: "RegisterUserPasswordsDontMatch",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			data: dto.RegisterRequest{
				Username:        "test",
				Email:           "test@ocg.com",
				Password:        "test",
				ConfirmPassword: "test123",
			},
			UserRepoResp:    nil,
			SessionServResp: "",
			SessionServErr:  nil,
			respError:       usecase_errors.BadRequestError{},
			mustErr:         true,
		},
		{
			testName: "RegisterUserEmailAlreadyExists",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			data: dto.RegisterRequest{
				Username:        "test",
				Email:           "test@ocg.com",
				Password:        "test",
				ConfirmPassword: "test",
			},
			UserRepoResp:    repositories.ErrDuplicate,
			SessionServResp: "",
			SessionServErr:  nil,
			respError:       usecase_errors.AlreadyExistsError{},
			mustErr:         true,
		},
		{
			testName: "RegisterUserSuccess",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			data: dto.RegisterRequest{
				Username:        "test",
				Email:           "test@ocg.com",
				Password:        "test",
				ConfirmPassword: "test",
			},
			UserRepoResp:    nil,
			SessionServResp: "test",
			SessionServErr:  nil,
			respError:       nil,
			mustErr:         false,
		},
	}

	for _, tc := range testCases {
		mockUserRepository := new(mocks.IUserRepository)
		mockSessionService := new(mocks.ISessionService)
		mockEmailService := new(mocks.IEmailService)
		service.UserRepository = mockUserRepository
		service.SessionService = mockSessionService
		service.EmailService = mockEmailService

		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository.EXPECT().Create(context.Background(), mock.Anything).Return(tc.UserRepoResp)
			mockSessionService.EXPECT().SetSession(context.Background(), mock.Anything).Return(tc.SessionServResp, tc.SessionServErr)
			mockEmailService.EXPECT().SendRegisterEmail(mock.Anything, mock.Anything).Maybe().Return(tc.respError)

			err := service.RegisterUser(context.Background(), tc.caller, tc.data)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.respError), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestLogin(t *testing.T) {
	mockApp := GetAppMock()
	service := services.AuthService{
		App: mockApp,
	}

	testCases := []struct {
		testName      string
		caller        dto.UserDTO
		data          dto.LoginRequest
		userRepoResp  []domain.User
		userRepoErr   error
		sessServResp  string
		sessServErr   error
		expectedResp  string
		expectedError error
		mustErr       bool
	}{
		{
			testName: "LoginUserAlreadyLoggedIn",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			expectedError: usecase_errors.BadRequestError{},
			mustErr:       true,
		},
		{
			testName: "LoginUserNotFound",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			data: dto.LoginRequest{
				UsernameOrEmail: "test",
				Password:        "test",
			},
			userRepoResp:  []domain.User{},
			expectedError: usecase_errors.BadRequestError{},
			mustErr:       true,
		},
		{
			testName: "LoginUserDbError",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			data: dto.LoginRequest{
				UsernameOrEmail: "test",
				Password:        "test",
			},
			userRepoResp:  []domain.User{},
			userRepoErr:   errors.New("internal db error"),
			expectedError: usecase_errors.BadRequestError{},
			mustErr:       true,
		},
		{
			testName: "LoginUserInvalidPassword",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			data: dto.LoginRequest{
				UsernameOrEmail: "test",
				Password:        "invalidPassword",
			},
			userRepoResp: []domain.User{
				{
					Username: "test",
					Password: func() string {
						hash, err := utils.HashPassword("test123")
						if err != nil {
							panic(err)
						}
						return hash
					}(),
					IsActive: true,
					Role:     enums.USER,
				},
			},
			expectedError: usecase_errors.BadRequestError{},
			mustErr:       true,
		},
		{
			testName: "LoginUserNotActive",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			data: dto.LoginRequest{
				UsernameOrEmail: "test",
				Password: func() string {
					hash, err := utils.HashPassword("test")
					if err != nil {
						panic(err)
					}
					return hash
				}(),
			},
			userRepoResp: []domain.User{
				{
					Username: "test",
					Password: func() string {
						hash, err := utils.HashPassword("test")
						if err != nil {
							panic(err)
						}
						return hash
					}(),
					IsActive: false,
					Role:     enums.ANONYMOUS,
				},
			},
			expectedError: usecase_errors.BadRequestError{},
			mustErr:       true,
		},
		{
			testName: "LoginUserSuccess",
			caller: dto.UserDTO{
				ID:       1,
				Role:     enums.ANONYMOUS,
				IsActive: false,
			},
			data: dto.LoginRequest{
				UsernameOrEmail: "test",
				Password:        "test",
			},
			userRepoResp: []domain.User{
				{
					Username: "test",
					Email:    "test@ocg.com",
					Password: func() string {
						hash, err := utils.HashPassword("test")
						if err != nil {
							panic(err)
						}
						return hash
					}(),
					IsActive: true,
					Role:     enums.USER,
				},
			},
			sessServResp: "test",
			expectedResp: "test",
			mustErr:      false,
		},
	}

	for _, tc := range testCases {
		mockUserRepository := new(mocks.IUserRepository)
		mockSessionService := new(mocks.ISessionService)
		service.UserRepository = mockUserRepository
		service.SessionService = mockSessionService

		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository.EXPECT().Filter(context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(tc.userRepoResp, tc.userRepoErr)
			mockSessionService.EXPECT().SetSession(context.Background(), mock.Anything).Return(tc.sessServResp, tc.sessServErr)

			res, err := service.Login(context.Background(), tc.caller, tc.data)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectedError), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResp, res)
			}
		})
	}
}
