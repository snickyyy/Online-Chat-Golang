package handler_api

import (
	"libs/src/internal/dto"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	bla, _ := ctx.Get("user")
	user, _ := bla.(dto.UserDTO)
	ctx.JSON(200, user)
}
