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
	service := services.NewUserService(settings.AppVar)

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
	sessId, err := c.Cookie("sessionID")
	var requestData dto.ChangeUserProfileRequest
	if err != nil {
		c.Error(api_errors.ErrNeedLoginForChangeProfile)
		return
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.Error(api_errors.ErrInvalidBody)
		return
	}

	service := services.NewUserService(settings.AppVar)
	err = service.ChangeUserProfile(requestData, sessId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ChangeUserProfileResponse{ChangedFields: requestData, Message: "success"})

}
