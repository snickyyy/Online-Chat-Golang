package handler_api

import (
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
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
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(200, profile)
}
