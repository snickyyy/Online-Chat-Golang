package handler_middlewares

import "github.com/gin-gonic/gin"


func HuiMiddleware(ctx *gin.Context) {
	id := ctx.Param("id")
	println("test1", id)
	ctx.Next()
	println("test2", id)
}