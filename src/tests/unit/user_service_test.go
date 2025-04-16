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
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/pkg/utils"
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
			expectedErr:  api_errors.ErrProfileNotFound,
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
			expectedErr: api_errors.ErrProfileNotFound,
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
			expectedErr: api_errors.ErrProfileNotFound,
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
			mockUserRepository.EXPECT().Filter(mock.Anything, mock.Anything).Return(tc.UserRepoResp, tc.UserRepoErr)

			res, err := service.GetUserProfile(tc.username)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
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
			expectedErr: api_errors.ErrNeedLoginForChangeProfile,
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
			expectedErr: api_errors.ErrNeedLoginForChangeProfile,
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
			expectedErr: api_errors.ErrUserAlreadyExists,
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
			mockUserRepository.EXPECT().UpdateById(mock.Anything, mock.Anything).Return(tc.userRepoErr)

			err := service.ChangeUserProfile(tc.caller, tc.data)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
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
			expectedErr:  api_errors.ErrUserNotFound,
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
			expectedErr: api_errors.ErrUserNotFound,
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
			expectedErr: api_errors.ErrUserNotFound,
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
			mockUserRepository.EXPECT().Filter(mock.Anything, mock.Anything, mock.Anything).Return(tc.userRepoResp, tc.userRepoErr)
			mockSessionService.EXPECT().SetSession(mock.Anything).Return(tc.sessionServResp, tc.sessionServErr)
			mockEmailService.EXPECT().SendResetPasswordEmail(mock.Anything, mock.Anything).Maybe().Return(nil)

			code, err := service.ResetPassword(tc.request)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
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
			expectedErr: api_errors.ErrUnauthorized,
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
			expectedErr: api_errors.ErrUnauthorized,
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
			expectedErr: api_errors.ErrInvalidPassword,
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
			expectedErr: api_errors.ErrSamePassword,
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
			expectedErr: api_errors.ErrPasswordsDontMatch,
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
			mockUserRepository.EXPECT().GetById(mock.Anything).Return(tc.userRepoGetResp, tc.userRepoGetErr)
			mockUserRepository.EXPECT().UpdateById(mock.Anything, mock.Anything).Return(tc.userRepoUpdateResp)

			err := service.ChangePassword(tc.caller, tc.request)

			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
