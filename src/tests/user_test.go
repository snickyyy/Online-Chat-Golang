package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	"libs/src/settings"
	"net/http"
)

func (suite *AppTestSuite) TestGetProfile() {
	url := "http://127.0.0.1:8000/accounts/profile/"

	userService := services.NewUserService(settings.AppVar)
	authService := services.NewAuthService(settings.AppVar)

	err := userService.CreateSuperUser("testprofile1", "profiletest1@test.com", "test123")
	suite.NoError(err)

	err = authService.RegisterUser(
		dto.RegisterRequest{
			Username:        "profileNonConfirm",
			Email:           "profileNonConfirm@test.com",
			Password:        "test123",
			ConfirmPassword: "test123",
		},
	)
	suite.NoError(err)

	res, err := suite.client.Get(url + "testprofile1")
	suite.NoError(err)
	suite.Equal(http.StatusOK, res.StatusCode)

	res, err = suite.client.Get(url + "testprofile1")
	suite.NoError(err)
	suite.Equal(http.StatusOK, res.StatusCode)

	res, err = suite.client.Get(url + "profileNonConfirm")
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, res.StatusCode)

	defer res.Body.Close()
}

func (suite *AppTestSuite) TestChangeProfile() {
	urlGetProfile := "http://127.0.0.1:8000/accounts/profile/"
	url := "http://127.0.0.1:8000/accounts/profile/edit"

	userService := services.NewUserService(settings.AppVar)
	authService := services.NewAuthService(settings.AppVar)

	err := userService.CreateSuperUser("TestProfileEdit", "profileEditTest@test.com", "test123")
	suite.NoError(err)

	sess, err := authService.Login(dto.LoginRequest{UsernameOrEmail: "TestProfileEdit", Password: "test123"})
	suite.NoError(err)

	changeBody := dto.ChangeUserProfileRequest{
		NewUsername: new(string),
	}
	*changeBody.NewUsername = "TestProfileEditNewUsername"
	body, _ := json.Marshal(changeBody)

	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	suite.NoError(err)
	resProfile, err := suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, resProfile.StatusCode)

	request, err = http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	suite.NoError(err)
	request.AddCookie(&http.Cookie{Name: "sessionID", Value: sess})
	_, err = suite.client.Do(request)
	suite.NoError(err)

	resProfile, err = suite.client.Get(urlGetProfile + "TestProfileEditNewUsername")
	suite.NoError(err)
	suite.Equal(http.StatusOK, resProfile.StatusCode)
	encode, _ := io.ReadAll(resProfile.Body)
	var profile dto.UserProfile
	json.Unmarshal(encode, &profile)
	suite.Equal(profile.Username, "TestProfileEditNewUsername")

	defer resProfile.Body.Close()
}
