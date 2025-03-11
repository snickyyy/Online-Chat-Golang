package tests

import (
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
			Username: "profileNonConfirm",
			Email:    "profileNonConfirm@test.com",
            Password: "test123",
            ConfirmPassword: "test123",
		},
	)
	suite.NoError(err)

	res, err := suite.client.Get(url + "testprofile1")
	suite.NoError(err)
	suite.Equal(http.StatusOK, res.StatusCode)

	res, err = suite.client.Get(url + "testprofile1")
	suite.NoError(err)
	suite.Equal(http.StatusOK, res.StatusCode) // blablabla
}
