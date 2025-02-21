package handler_middlewares

import (
	"libs/src/settings"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware(ctx *gin.Context) {
	settings.AppVar.Logger.Info("Test1")
	ctx.Next()
	settings.AppVar.Logger.Info("Test2")
}