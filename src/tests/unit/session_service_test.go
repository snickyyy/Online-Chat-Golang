package unit

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"libs/src/internal/dto"
	"libs/src/internal/mocks"
	services "libs/src/internal/usecase"
	api_errors "libs/src/internal/usecase/errors"
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
