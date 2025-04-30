package handler_api

import (
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	"libs/src/settings"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	bla, _ := ctx.Get("user")
	user, _ := bla.(dto.UserDTO)
	ctx.JSON(200, user)
}

// @Summary Ping
// @Description setting user online
// @Tags Base
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ping [get]
func Ping(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	userService := services.NewUserService(app)
	if err := userService.SetOnline(c.Request.Context(), user); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, dto.MessageResponse{Message: "pong"})
}
