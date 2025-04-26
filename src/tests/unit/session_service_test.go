package unit

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"libs/src/internal/domain/enums"
	"libs/src/internal/dto"
	"libs/src/internal/mocks"
	services "libs/src/internal/usecase"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/pkg/utils"
	"reflect"
	"testing"
)

func TestGetSession(t *testing.T) {
	mockApp := GetAppMock()

	sessionService := services.SessionService{
		App: mockApp,
	}

	testCases := []struct {
		testName    string
		Prefix      string
		Session     string
		mockResp    string
		mockErr     error
		mustResp    dto.SessionDTO
		mustRespErr error
		mustErr     bool
	}{
		{"TestGetSessionInvalidSession", "test", "test", "", errors.New("Invalid session"), dto.SessionDTO{}, usecase_errors.BadRequestError{}, true},
		{"TestGetSessionValidSession", "test", "test", `{"id":"test","prefix":"test","payload":"test"}`, nil, dto.SessionDTO{SessionID: "test", Prefix: "test"}, nil, false},
	}
	for _, tc := range testCases {
		mockRedisRepo := new(mocks.IBaseRedisRepository)
		sessionService.RedisBaseRepository = mockRedisRepo

		t.Run(tc.testName, func(t *testing.T) {
			mockRedisRepo.EXPECT().GetByKey(context.Background(), tc.Prefix, tc.Session).Return(tc.mockResp, tc.mockErr)

			resp, err := sessionService.GetSession(context.Background(), tc.Prefix, tc.Session)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.mustRespErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mustResp.SessionID, resp.SessionID)
				assert.Equal(t, tc.mustResp.Prefix, resp.Prefix)
			}

		})
	}
}
func TestDeleteSession(t *testing.T) {
	mockApp := GetAppMock()
	sessionService := services.SessionService{
		App: mockApp,
	}

	testCases := []struct {
		testName    string
		Prefix      string
		Session     string
		mockResp    int64
		mockErr     error
		mustRespErr error
		mustErr     bool
	}{
		{
			"TestDeleteSessionInvalidSession", "test", "test", 0, errors.New("Invalid session"), usecase_errors.BadRequestError{}, true,
		},
		{
			"TestDeleteSessionValidSession", "test", "test", 1, nil, nil, false,
		},
	}
	for _, tc := range testCases {
		mockRedisRepo := new(mocks.IBaseRedisRepository)
		sessionService.RedisBaseRepository = mockRedisRepo

		t.Run(tc.testName, func(t *testing.T) {
			mockRedisRepo.EXPECT().Delete(context.Background(), tc.Prefix, tc.Session).Return(tc.mockResp, tc.mockErr)

			err := sessionService.DeleteSession(context.Background(), tc.Prefix, tc.Session)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.mustRespErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestGetUserByAuthSession(t *testing.T) {
	mockApp := GetAppMock()
	sessionService := services.SessionService{
		App: mockApp,
	}
	testCases := []struct {
		testName    string
		Session     string
		mockResp    string
		mockErr     error
		mustResp    dto.UserDTO
		mustRespErr error
		mustErr     bool
	}{
		{
			"TestGetUserByAuthSessionInvalidSession",
			"test",
			"",
			errors.New("Invalid session"),
			dto.UserDTO{},
			usecase_errors.BadRequestError{},
			true,
		},
		{
			"TestGetUserByAuthSessionSuccess",
			"test",
			func() string {
				user := dto.UserDTO{
					ID:       1,
					Username: "test123",
					Email:    "test@test.ocg",
					IsActive: true,
					Role:     enums.USER,
				}
				encoding, _ := json.Marshal(
					dto.AuthSession{
						UserDTO: user,
					},
				)
				encrypt, err := utils.Encrypt(mockApp.Config.AppConfig.SecretKey, string(encoding))
				if err != nil {
					t.Fatal("Error encrypting auth session:", err)
				}
				session := dto.SessionDTO{
					SessionID: "test",
					Prefix:    "session:",
					Payload:   encrypt,
				}
				encode, _ := json.Marshal(session)
				return string(encode)
			}(),
			nil,
			dto.UserDTO{
				ID:       1,
				Username: "test123",
				Email:    "test@test.ocg",
				IsActive: true,
				Role:     enums.USER,
			},
			nil,
			false,
		},
	}
	for _, tc := range testCases {
		mockRepo := new(mocks.IBaseRedisRepository)
		sessionService.RedisBaseRepository = mockRepo

		t.Run(tc.testName, func(t *testing.T) {
			mockRepo.EXPECT().GetByKey(context.Background(), mockApp.Config.RedisConfig.Prefixes.SessionPrefix, tc.Session).Return(tc.mockResp, tc.mockErr)

			res, err := sessionService.GetUserByAuthSession(context.Background(), tc.Session)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.mustRespErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mustResp.ID, res.ID)
				assert.Equal(t, tc.mustResp.Username, res.Username)
				assert.Equal(t, tc.mustResp.Email, res.Email)
			}
		})
	}
}

func TestGetUserByEmailSession(t *testing.T) {
	mockApp := GetAppMock()
	sessionService := services.SessionService{
		App: mockApp,
	}
	testCases := []struct {
		testName    string
		Session     string
		mockResp    string
		mockErr     error
		mustResp    dto.UserDTO
		mustRespErr error
		mustErr     bool
	}{
		{
			"TestGetUserByAuthSessionInvalidSession",
			"test",
			"",
			errors.New("Invalid session"),
			dto.UserDTO{},
			usecase_errors.BadRequestError{},
			true,
		},
		{
			"TestGetUserByAuthSessionSuccess",
			"test",
			func() string {
				user := dto.UserDTO{
					ID:       1,
					Username: "test123",
					Email:    "test@test.ocg",
					IsActive: true,
					Role:     enums.USER,
				}
				encoding, _ := json.Marshal(
					dto.EmailSession{
						UserDTO: user,
					},
				)
				encrypt, err := utils.Encrypt(mockApp.Config.AppConfig.SecretKey, string(encoding))
				if err != nil {
					t.Fatal("Error encrypting auth session:", err)
				}
				session := dto.SessionDTO{
					SessionID: "test",
					Prefix:    "session:",
					Payload:   encrypt,
				}
				encode, _ := json.Marshal(session)
				return string(encode)
			}(),
			nil,
			dto.UserDTO{
				ID:       1,
				Username: "test123",
				Email:    "test@test.ocg",
				IsActive: true,
				Role:     enums.USER,
			},
			nil,
			false,
		},
	}
	for _, tc := range testCases {
		mockRepo := new(mocks.IBaseRedisRepository)
		sessionService.RedisBaseRepository = mockRepo

		t.Run(tc.testName, func(t *testing.T) {
			mockRepo.EXPECT().GetByKey(context.Background(), mockApp.Config.RedisConfig.Prefixes.SessionPrefix, tc.Session).Return(tc.mockResp, tc.mockErr)

			res, err := sessionService.GetUserByAuthSession(context.Background(), tc.Session)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, reflect.TypeOf(tc.mustRespErr), reflect.TypeOf(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mustResp.ID, res.ID)
				assert.Equal(t, tc.mustResp.Username, res.Username)
				assert.Equal(t, tc.mustResp.Email, res.Email)
			}
		})
	}
}
