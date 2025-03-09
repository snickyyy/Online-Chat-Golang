package handler_middlewares

import (
	"fmt"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
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
		c.Set("user.state.isActive", false)
		c.Next()
		return
	}

	service := services.NewAuthService(settings.AppVar)

	user, err = service.GetUserBySession(sid)
	if err != nil {
		settings.AppVar.Logger.Error(fmt.Sprintf("Error getting session: %s || error: %s", sid, err))
		user = unknown
		c.Set("user.state.isActive", false)
		c.Next()
		return
	}

	if !user.IsActive || user.Role == domain.ANONYMOUS {
		user = unknown
		c.Set("user.state.isActive", false)
		c.Next()
		return
	}
	c.Set("user.state.isActive", true)
	c.Set("user", user)
	c.Next()
}
