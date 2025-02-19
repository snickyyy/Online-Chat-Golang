package handler_api

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	time.Sleep(time.Second / 2)
	ctx.JSON(200, "Hello, Gin!")
}
