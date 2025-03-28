package tests

import (
	"bytes"
	"encoding/json"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	"libs/src/settings"
	"net/http"
)

func (suite *AppTestSuite) TestCreateChat() {
	url := "http://127.0.0.1:8000/messenger/chat/create"

	userService := services.NewUserService(settings.AppVar)
	authService := services.NewAuthService(settings.AppVar)

	err := userService.CreateSuperUser("TestCreateChat", "TestCreateChat@test.com", "test123")
	suite.NoError(err)
	sess, err := authService.Login(dto.LoginRequest{UsernameOrEmail: "TestCreateChat", Password: "test123"})
	suite.NoError(err)

	dataCreateChat, err := json.Marshal(dto.CreateChatRequest{
		Title:       "TestCreateChat",
		Description: "TestCreateChat",
	})
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(dataCreateChat))
	requestResult, err := suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, requestResult.StatusCode)

	request, _ = http.NewRequest("POST", url, bytes.NewBuffer(dataCreateChat))
	request.AddCookie(&http.Cookie{Name: "sessionID", Value: sess})
	requestResult, err = suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusOK, requestResult.StatusCode)
}
