package unit

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"libs/src/internal/domain/enums"
	"libs/src/internal/dto"
	"libs/src/internal/mocks"
	services "libs/src/internal/usecase"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/internal/usecase/utils"
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
		{"TestGetSessionInvalidSession", "test", "test", "", errors.New("Invalid session"), dto.SessionDTO{}, api_errors.ErrInvalidSession, true},
		{"TestGetSessionValidSession", "test", "test", `{"id":"test","prefix":"test","payload":"test"}`, nil, dto.SessionDTO{SessionID: "test", Prefix: "test"}, nil, false},
	}
	for _, tc := range testCases {
		mockRedisRepo := new(mocks.IBaseRedisRepository)
		sessionService.RedisBaseRepository = mockRedisRepo

		t.Run(tc.testName, func(t *testing.T) {
			mockRedisRepo.EXPECT().GetByKey(tc.Prefix, tc.Session).Return(tc.mockResp, tc.mockErr)

			resp, err := sessionService.GetSession(tc.Prefix, tc.Session)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.mustRespErr, err)
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
			"TestDeleteSessionInvalidSession", "test", "test", 0, errors.New("Invalid session"), api_errors.ErrInvalidSession, true,
		},
		{
			"TestDeleteSessionValidSession", "test", "test", 1, nil, nil, false,
		},
	}
	for _, tc := range testCases {
		mockRedisRepo := new(mocks.IBaseRedisRepository)
		sessionService.RedisBaseRepository = mockRedisRepo

		t.Run(tc.testName, func(t *testing.T) {
			mockRedisRepo.EXPECT().Delete(tc.Prefix, tc.Session).Return(tc.mockResp, tc.mockErr)

			err := sessionService.DeleteSession(tc.Prefix, tc.Session)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.mustRespErr, err)
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
			api_errors.ErrInvalidSession,
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
			mockRepo.EXPECT().GetByKey(mockApp.Config.RedisConfig.Prefixes.SessionPrefix, tc.Session).Return(tc.mockResp, tc.mockErr)

			res, err := sessionService.GetUserByAuthSession(tc.Session)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.mustRespErr, err)
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
			api_errors.ErrInvalidSession,
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
			mockRepo.EXPECT().GetByKey(mockApp.Config.RedisConfig.Prefixes.SessionPrefix, tc.Session).Return(tc.mockResp, tc.mockErr)

			res, err := sessionService.GetUserByAuthSession(tc.Session)
			if tc.mustErr {
				assert.Error(t, err)
				assert.Equal(t, tc.mustRespErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mustResp.ID, res.ID)
				assert.Equal(t, tc.mustResp.Username, res.Username)
				assert.Equal(t, tc.mustResp.Email, res.Email)
			}
		})
	}
}
