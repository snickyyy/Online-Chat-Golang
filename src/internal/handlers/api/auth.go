package handler_api

import (
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	services "libs/src/internal/usecase"
	"libs/src/settings"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	
	service := services.AuthSessionService{
		RedisBaseRepository: repositories.BaseRedisRepository{
			Client: settings.AppVar.RedisSess,
			Ctx:    settings.Context.Ctx,
		},
		App: settings.AppVar,
	}
	id, err := service.SetSession(dto.UserDTO{
		UserName: "snicky",
		Email:    "snicky@blabla.com",
	},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"session_id": id})
}
