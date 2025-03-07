package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"libs/src/internal/dto"
	api_errors "libs/src/internal/usecase/errors"
	"log"
	"net/http"
)

func (suite *AppTestSuite) TestRegister() {
	log.Println("Testing Register endpoint...")
	url := "http://127.0.0.1:8000/accounts/auth/register"
	contType := "application/json"

	dataSuccess, _ := json.Marshal(dto.RegisterRequest{
		Username: "test",
		Email:    "test@example.com",
		Password: "test",
		ConfirmPassword: "test",
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
		Username: "test",
		Email:    "test@example.com",
		Password: "test",
		ConfirmPassword: "test",
	})
	res, err = suite.client.Post(url, contType, bytes.NewBuffer(dataUnique))
	suite.NoError(err)
	responseBody, _ := ioutil.ReadAll(res.Body)
	suite.NoError(err)
	suite.Equal(string(responseBody), fmt.Sprintf(`{"error":"%s"}`, api_errors.ErrUserAlreadyExists))
	suite.Equal(http.StatusBadRequest, res.StatusCode)
	log.Println("Response status code: ", res.StatusCode)

}
