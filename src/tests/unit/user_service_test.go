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
	"libs/src/pkg/utils"
	"reflect"
	"testing"
)

func TestGetUserProfile(t *testing.T) {
	mockApp := GetAppMock()
	service := services.UserService{
		App: mockApp,
	}

	testCases := []struct {
		testName     string
		username     string
		UserRepoResp []domain.User
		UserRepoErr  error
		expectedResp dto.UserProfile
		expectedErr  error
		mustErr      bool
	}{
		{
			testName:     "GetUserProfileNotFound",
			username:     "testuser",
			UserRepoResp: []domain.User{},
			expectedErr:  usecase_errors.NotFoundError{},
			mustErr:      true,
		},
		{
			testName: "GetUserProfileNotAuthenticated",
			username: "testuser",
			UserRepoResp: []domain.User{
				{
					Username: "testuser",
					IsActive: true,
					Role:     enums.ANONYMOUS,
				},
			},
			expectedErr: usecase_errors.NotFoundError{},
			mustErr:     true,
		},
		{
			testName: "GetUserProfileNotActive",
			username: "testuser",
			UserRepoResp: []domain.User{
				{
					Username: "testuser",
					IsActive: false,
					Role:     enums.USER,
				},
			},
			expectedErr: usecase_errors.NotFoundError{},
			mustErr:     true,
		},
		{
			testName: "GetUserProfileSuccess",
			username: "testuser",
			UserRepoResp: []domain.User{
				{
					Username:    "testuser",
					Email:       "testuser@example.com",
					IsActive:    true,
					Role:        enums.USER,
					Description: "test description",
				},
			},
			expectedResp: dto.UserProfile{
				Username:    "testuser",
				Description: "test description",
				Role:        enums.RolesToLabels[enums.USER],
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockUserRepository := new(mocks.IUserRepository)
		service.UserRepository = mockUserRepository

		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository.EXPECT().Filter(mockApp.Ctx, mock.Anything, mock.Anything).Return(tc.UserRepoResp, tc.UserRepoErr)

			res, err := service.GetUserProfile(mockApp.Ctx, tc.username)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectedErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResp.Username, res.Username)
				assert.Equal(t, tc.expectedResp.Description, res.Description)
				assert.Equal(t, tc.expectedResp.Role, res.Role)
			}
		})
	}
}
func TestChangeUserProfile(t *testing.T) {
	mockApp := GetAppMock()
	service := services.UserService{
		App: mockApp,
	}

	username := "testuser"
	description := "test description"
	data := dto.ChangeUserProfileRequest{
		NewUsername:    &username,
		NewDescription: &description,
	}

	testCases := []struct {
		testName    string
		data        dto.ChangeUserProfileRequest
		caller      dto.UserDTO
		userRepoErr error
		expectedErr error
		mustErr     bool
	}{
		{
			testName: "ChangeUserProfileNotAuthenticated",
			data:     data,
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			expectedErr: usecase_errors.UnauthorizedError{},
			mustErr:     true,
		},
		{
			testName: "ChangeUserProfileNotActive",
			data:     data,
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: false,
			},
			expectedErr: usecase_errors.UnauthorizedError{},
			mustErr:     true,
		},
		{
			testName: "ChangeUserProfileDbErrDuplicateUsername",
			data:     data,
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: true,
				Email:    "testuser@example.com",
			},
			userRepoErr: repositories.ErrDuplicate,
			expectedErr: usecase_errors.AlreadyExistsError{},
			mustErr:     true,
		},
		{
			testName: "ChangeUserProfileSuccess",
			data:     data,
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: true,
				Email:    "testuser@example.com",
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockUserRepository := new(mocks.IUserRepository)

		service.UserRepository = mockUserRepository

		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository.EXPECT().UpdateById(mockApp.Ctx, mock.Anything, mock.Anything).Return(tc.userRepoErr)

			err := service.ChangeUserProfile(mockApp.Ctx, tc.caller, tc.data)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectedErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestResetPassword(t *testing.T) {
	mockApp := GetAppMock()
	service := services.UserService{
		App: mockApp,
	}

	request := dto.ResetPasswordRequest{
		UsernameOrEmail: "testuser",
	}

	sessServErr := errors.New("internal session error")

	testCases := []struct {
		testName        string
		request         dto.ResetPasswordRequest
		userRepoResp    []domain.User
		userRepoErr     error
		sessionServResp string
		sessionServErr  error
		expectedErr     error
		mustErr         bool
	}{
		{
			testName:     "ResetPasswordUserNotFound",
			request:      request,
			userRepoResp: []domain.User{},
			expectedErr:  usecase_errors.NotFoundError{},
			mustErr:      true,
		},
		{
			testName: "ResetPasswordUserNotActive",
			request:  request,
			userRepoResp: []domain.User{
				{
					Username: "testuser",
					Role:     enums.USER,
					IsActive: false,
				},
			},
			expectedErr: usecase_errors.NotFoundError{},
			mustErr:     true,
		},
		{
			testName: "ResetPasswordUserNotAuthenticated",
			request:  request,
			userRepoResp: []domain.User{
				{
					Username: "testuser",
					Role:     enums.ANONYMOUS,
					IsActive: true,
				},
			},
			expectedErr: usecase_errors.NotFoundError{},
			mustErr:     true,
		},
		{
			testName: "ResetPasswordSessionError",
			request:  request,
			userRepoResp: []domain.User{
				{
					Username: "testuser",
					Role:     enums.USER,
					IsActive: true,
				},
			},
			sessionServErr: sessServErr,
			expectedErr:    sessServErr,
			mustErr:        true,
		},
		{
			testName: "ResetPasswordSuccess",
			request:  request,
			userRepoResp: []domain.User{
				{
					Username: "testuser",
					Role:     enums.USER,
					IsActive: true,
				},
			},
			sessionServResp: "sessionId",
			mustErr:         false,
		},
	}
	for _, tc := range testCases {
		mockSessionService := new(mocks.ISessionService)
		mockUserRepository := new(mocks.IUserRepository)
		mockEmailService := new(mocks.IEmailService)
		service.SessionService = mockSessionService
		service.UserRepository = mockUserRepository
		service.EmailService = mockEmailService

		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository.EXPECT().Filter(mockApp.Ctx, mock.Anything, mock.Anything, mock.Anything).Return(tc.userRepoResp, tc.userRepoErr)
			mockSessionService.EXPECT().SetSession(mockApp.Ctx, mock.Anything).Return(tc.sessionServResp, tc.sessionServErr)
			mockEmailService.EXPECT().SendResetPasswordEmail(mock.Anything, mock.Anything).Maybe().Return(nil)

			code, err := service.ResetPassword(mockApp.Ctx, tc.request)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectedErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, 0, code)
			}
		})
	}
}
func TestChangePassword(t *testing.T) {
	mockApp := GetAppMock()
	service := services.UserService{
		App: mockApp,
	}

	request := dto.ChangePasswordRequest{
		OldPassword:        "test123",
		NewPassword:        "test",
		ConfirmNewPassword: "test",
	}

	testCases := []struct {
		testName string

		caller  dto.UserDTO
		request dto.ChangePasswordRequest

		userRepoGetResp domain.User
		userRepoGetErr  error

		userRepoUpdateResp error

		expectedErr error
		mustErr     bool
	}{
		{
			testName: "ChangePasswordNotAuthenticated",
			request:  request,
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			userRepoGetResp: domain.User{
				Username: "testuser",
				Role:     enums.ANONYMOUS,
				IsActive: true,
			},
			expectedErr: usecase_errors.UnauthorizedError{},
			mustErr:     true,
		},
		{
			testName: "ChangePasswordNotActive",
			request:  request,
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: false,
			},
			userRepoGetResp: domain.User{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: false,
			},
			expectedErr: usecase_errors.UnauthorizedError{},
			mustErr:     true,
		},
		{
			testName: "ChangePasswordInvalidOldPassword",
			request:  request,
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: true,
			},
			userRepoGetResp: domain.User{
				Username: "testuser",
				Email:    "test@ocg.com",
				Role:     enums.USER,
				IsActive: true,
				Password: func() string {
					pass, _ := utils.HashPassword("invalid")
					return pass
				}(),
			},
			expectedErr: usecase_errors.BadRequestError{},
			mustErr:     true,
		},
		{
			testName: "ChangePasswordSamePassword",
			request: dto.ChangePasswordRequest{
				OldPassword:        "test123",
				NewPassword:        "test123",
				ConfirmNewPassword: "test123",
			},
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: true,
			},
			userRepoGetResp: domain.User{
				Username: "testuser",
				Email:    "test@ocg.com",
				Role:     enums.USER,
				IsActive: true,
				Password: func() string {
					pass, _ := utils.HashPassword("test123")
					return pass
				}(),
			},
			expectedErr: usecase_errors.BadRequestError{},
			mustErr:     true,
		},
		{
			testName: "ChangePasswordPasswordsDontMatch",
			request: dto.ChangePasswordRequest{
				OldPassword:        "test123",
				NewPassword:        "test",
				ConfirmNewPassword: "test123",
			},
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: true,
			},
			userRepoGetResp: domain.User{
				Username: "testuser",
				Email:    "test@ocg.com",
				Role:     enums.USER,
				IsActive: true,
				Password: func() string {
					pass, _ := utils.HashPassword("test123")
					return pass
				}(),
			},
			expectedErr: usecase_errors.BadRequestError{},
			mustErr:     true,
		},
		{
			testName: "ChangePasswordSuccess",
			request:  request,
			caller: dto.UserDTO{
				Username: "testuser",
				Role:     enums.USER,
				IsActive: true,
			},
			userRepoGetResp: domain.User{
				Username: "testuser",
				Email:    "test@ocg.com",
				Role:     enums.USER,
				IsActive: true,
				Password: func() string {
					pass, _ := utils.HashPassword("test123")
					return pass
				}(),
			},
			mustErr: false,
		},
	}

	for _, tc := range testCases {
		mockUserRepository := new(mocks.IUserRepository)

		service.UserRepository = mockUserRepository

		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository.EXPECT().GetById(mockApp.Ctx, mock.Anything).Return(tc.userRepoGetResp, tc.userRepoGetErr)
			mockUserRepository.EXPECT().UpdateById(mockApp.Ctx, mock.Anything, mock.Anything).Return(tc.userRepoUpdateResp)

			err := service.ChangePassword(mockApp.Ctx, tc.caller, tc.request)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expectedErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
