package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	"libs/src/settings"
	"mime/multipart"
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

	request, err := http.NewRequest("PATCH", url, &bytes.Buffer{})
	suite.NoError(err)
	resProfile, err := suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, resProfile.StatusCode)

	var changeBody bytes.Buffer
	writer := multipart.NewWriter(&changeBody)
	writer.WriteField("new_username", "TestProfileEditNewUsername")
	err = writer.Close()
	suite.NoError(err)
	request, err = http.NewRequest("PATCH", url, &changeBody)
	suite.NoError(err)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.NoError(err)
	request.AddCookie(&http.Cookie{Name: "sessionID", Value: sess})
	resChangeProfile, err := suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resChangeProfile.StatusCode)
	bla, _ := io.ReadAll(resChangeProfile.Body)
	fmt.Println(string(bla))

	resProfile, err = suite.client.Get(urlGetProfile + "TestProfileEditNewUsername")
	suite.NoError(err)
	suite.Equal(http.StatusOK, resProfile.StatusCode)
	encode, _ := io.ReadAll(resProfile.Body)
	var profile dto.UserProfile
	json.Unmarshal(encode, &profile)
	suite.Equal(profile.Username, "TestProfileEditNewUsername")

	defer resProfile.Body.Close()
}

func (suite *AppTestSuite) TestResetPassword() {
	// TEST RESET PASSWORD
	url := "http://127.0.0.1:8000/accounts/profile/reset-password"

	userService := services.NewUserService(settings.AppVar)

	err := userService.CreateSuperUser("TestResetPassword", "TestResetPassword@test.com", "test123")
	suite.NoError(err)

	changeBody := dto.ResetPasswordRequest{
		UsernameOrEmail: "TestResetPassword",
	}
	body, _ := json.Marshal(changeBody)

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	suite.NoError(err)
	resetPassResult, err := suite.client.Do(request)
	suite.NoError(err)

	suite.Equal(http.StatusOK, resetPassResult.StatusCode)

	encode, _ := io.ReadAll(resetPassResult.Body)
	var resetPassResponse dto.ResetPasswordResponse // SUCCESSFUL DATA
	json.Unmarshal(encode, &resetPassResponse)

	// TEST RESET PASSWORD WITH INVALID USERNAME
	changeBody = dto.ResetPasswordRequest{
		UsernameOrEmail: "unknownResetPassword",
	}
	body, _ = json.Marshal(changeBody)
	request, err = http.NewRequest("PUT", url, bytes.NewBuffer(body))
	suite.NoError(err)
	badResetPassResult, err := suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, badResetPassResult.StatusCode)
}
