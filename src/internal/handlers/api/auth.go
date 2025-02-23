package handler_api

import (
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	services "libs/src/internal/usecase"
	"libs/src/settings"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var registerData dto.RegisterRequest

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
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
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}
