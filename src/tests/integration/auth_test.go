package integration

import (
	"bytes"
	"encoding/json"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	"libs/src/settings"
	"log"
	"net/http"
)

func (suite *AppTestSuite) TestRegister() {
	log.Println("Testing Register endpoint...")
	url := "http://127.0.0.1:8000/accounts/auth/register"
	contType := "application/json"

	dataSuccess, _ := json.Marshal(dto.RegisterRequest{
		Username:        "test",
		Email:           "test@example.com",
		Password:        "Test@123.",
		ConfirmPassword: "Test@123.",
	})

	res, err := suite.client.Post(url, contType, bytes.NewBuffer(dataSuccess))
	suite.NoError(err)
	suite.Equal(http.StatusOK, res.StatusCode)
	log.Println("Response status code: ", res.StatusCode)

	dataFail, _ := json.Marshal(dto.RegisterRequest{})
	res, err = suite.client.Post(url, contType, bytes.NewBuffer(dataFail))
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, res.StatusCode)
	log.Println("Response status code: ", res.StatusCode)

	dataUnique, _ := json.Marshal(dto.RegisterRequest{
		Username:        "test",
		Email:           "test@example.com",
		Password:        "Test@123.",
		ConfirmPassword: "Test@123.",
	})
	res, err = suite.client.Post(url, contType, bytes.NewBuffer(dataUnique))
	suite.NoError(err)
	suite.NoError(err)
	suite.Equal(http.StatusConflict, res.StatusCode)
	log.Println("Response status code: ", res.StatusCode)

	defer res.Body.Close()

}

func (suite *AppTestSuite) TestLogin() {
	log.Println("Testing Login endpoint...")
	url := "http://127.0.0.1:8000/accounts/auth/login"
	contType := "application/json"

	err := services.NewUserService(settings.AppVar).CreateSuperUser(suite.Ctx, "testuser", "test@test.com", "test123")
	if err != nil {
		log.Fatal(err)
	}

	dataSuccess, _ := json.Marshal(dto.LoginRequest{
		UsernameOrEmail: "testuser",
		Password:        "test123",
	})

	res, err := suite.client.Post(url, contType, bytes.NewBuffer(dataSuccess))
	suite.NoError(err)
	suite.Equal(http.StatusOK, res.StatusCode)
	cookieFromLogin := res.Cookies()
	suite.NotEmpty(cookieFromLogin)
	log.Println("Response status code: ", res.StatusCode, res.Cookies())

	dataFail, _ := json.Marshal(dto.LoginRequest{
		UsernameOrEmail: "testuser",
		Password:        "false_password",
	})
	res, err = suite.client.Post(url, contType, bytes.NewBuffer(dataFail))
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, res.StatusCode)

	// TEST LOGOUT
	urlLogout := "http://127.0.0.1:8000/accounts/auth/logout"
	req, err := http.NewRequest("DELETE", urlLogout, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.AddCookie(cookieFromLogin[0])

	resLogout, err := suite.client.Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resLogout.StatusCode)
	log.Println("Response status code: ", resLogout.StatusCode)

	defer res.Body.Close()
}
