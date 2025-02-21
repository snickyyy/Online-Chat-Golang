package handler_api

import (
	"libs/src/internal/repositories"
	"libs/src/settings"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	repo := repositories.BaseRedisRepository{
		Client: settings.AppVar.RedisSess,
	}
	res, err := repo.Create("fack", "chack")
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error() + "blabla"})
		return
	}
	ctx.JSON(200, res)
}
