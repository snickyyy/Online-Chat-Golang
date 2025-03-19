package handler_middlewares

import (
	"fmt"
	"libs/src/internal/domain/enums"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	"libs/src/settings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)

	var user dto.UserDTO

	unknown := dto.UserDTO{
		ID:       0,
		Role:     enums.ANONYMOUS,
		IsActive: false,
	}

	sid, err := c.Cookie("sessionID")
	if err != nil {
		user = unknown
		c.Set("user", user)
		c.Set("user.state.isActive", false)
		c.Next()
		return
	}

	service := services.NewAuthService(app)

	user, err = service.GetUserBySession(app.Config.RedisConfig.Prefixes.SessionPrefix, sid)
	if err != nil {
		app.Logger.Error(fmt.Sprintf("Error getting session: %s || error: %s", app.Config.RedisConfig.Prefixes.SessionPrefix+sid, err))
		user = unknown
		c.Set("user.state.isActive", false)
		c.Next()
		return
	}

	if !user.IsActive || user.Role == enums.ANONYMOUS {
		user = unknown
		c.Set("user.state.isActive", false)
		c.Next()
		return
	}
	c.Set("user.state.isActive", true)
	c.Set("user", user)
	c.Next()
}
