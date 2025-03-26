package handler_middlewares

import (
	"github.com/gin-gonic/gin"
	"libs/src/settings"
)

func DependenciesMiddleware(c *gin.Context) {
	c.Set("app", settings.AppVar)
	c.Next()
}
