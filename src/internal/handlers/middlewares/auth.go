package handler_middlewares

import (
	"fmt"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	services "libs/src/internal/usecase"
	"libs/src/settings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	var user dto.UserDTO

	unknown := dto.UserDTO{
		ID:       0,
		Role:     domain.ANONYMOUS,
		IsActive: false,
	}

	sid, err := c.Cookie("sessionID")
	if err != nil {
		user = unknown
		c.Set("user", user)
		c.Next()
		return
	}

	service := services.AuthService{
		RedisBaseRepository: repositories.BaseRedisRepository{
			Client: settings.AppVar.RedisSess,
			Ctx:    settings.Context.Ctx,
		},
		App: settings.AppVar,
	}

	user, err = service.GetUserBySession(sid)
	if err != nil {
		settings.AppVar.Logger.Error(fmt.Sprintf("Error getting session: %s", sid))
	}

	if !user.IsActive || user.Role == domain.ANONYMOUS {
		user = unknown
		c.Set("user", user)
		c.Next()
		return
	}
	c.Set("user", user)
	c.Next()
}
