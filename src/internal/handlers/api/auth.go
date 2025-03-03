package handler_api

import (
	_ "libs/src/docs"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	services "libs/src/internal/usecase"
	"libs/src/settings"

	"github.com/gin-gonic/gin"
)

// @Summary User register
// @Description Register a new user behind the selected fields
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "Data to register"
// @Success 200 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/auth/register [post]
func Register(c *gin.Context) {
	var registerData dto.RegisterRequest

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return
	}

	service := services.AuthService{
		RedisBaseRepository: repositories.BaseRedisRepository{
			Client: settings.AppVar.RedisSess,
			Ctx:    settings.Context.Ctx,
		},
		App: settings.AppVar,
	}

	err := service.RegisterUser(registerData)
	if err != nil {
		c.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, dto.RegisterResponse{Message: "success", Status: true})
}

// @Summary User confirm registration
// @Description Confirm users email
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "Data"
// @Success 200 {object} string "success"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/auth/confirm-account [get]
func ConfirmAccount(c *gin.Context) {
	session_id := c.Param("token")

	service := services.AuthService{
		RedisBaseRepository: repositories.BaseRedisRepository{
			Client: settings.AppVar.RedisSess,
			Ctx:    settings.Context.Ctx,
		},
		App: settings.AppVar,
	}
	sess, err := service.ConfirmAccount(session_id)
	if err != nil {
		c.JSON(409, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.SetCookie("sessionID", sess, int(settings.AppVar.Config.AuthConfig.AuthSessionTTL), "/", "", true, true)
	c.JSON(200, "success")
}
