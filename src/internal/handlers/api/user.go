package handler_api

import (
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Profile
// @Description View user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param username path string true "Username of the user"
// @Success 200 {object} dto.UserProfile
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/profile/{username} [get]
func UserProfile(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	service := services.NewUserService(app)

	username := c.Param("username")

	profile, err := service.GetUserProfile(username)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, profile)
}

// @Summary Edit profile
// @Description Edit user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param user body dto.ChangeUserProfileRequest true "Data"
// @Success 200 {object} dto.UserProfile
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/profile/edit [patch]
func ChangeUserProfile(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	sessId, err := c.Cookie("sessionID")
	var requestData dto.ChangeUserProfileRequest
	if err != nil {
		c.Error(api_errors.ErrNeedLoginForChangeProfile)
		return
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewUserService(app)
	err = service.ChangeUserProfile(requestData, sessId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ChangeUserProfileResponse{ChangedFields: requestData, Message: "success"})
}

// @Summary Reset password
// @Description Reset user password
// @Tags profile
// @Accept json
// @Produce json
// @Param user body dto.ResetPasswordRequest true "Data"
// @Success 200 {object} dto.ResetPasswordResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/profile/reset-password [put]
func ResetPassword(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	var requestData dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewUserService(app)
	code, err := service.ResetPassword(requestData)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ResetPasswordResponse{Message: "success", Code: code})
}

// @Summary confirm reset password
// @Description Confirm reset user password
// @Tags profile
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Param user body dto.ConfirmResetPasswordRequest true "Data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/profile/reset-password/confirm/{token} [put]
func ConfirmResetPassword(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	token := c.Param("token")
	var requestData dto.ConfirmResetPasswordRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewUserService(app)
	err := service.ConfirmResetPassword(token, requestData)
	if err != nil {
		c.Error(err)
		return
	}
	c.SetCookie("sessionID", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}
