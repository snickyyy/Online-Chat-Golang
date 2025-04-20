package handler_middlewares

import (
	"github.com/gin-gonic/gin"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
	"net/http"
)

func ErrorHandler(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	c.Next()
	err := c.Errors.Last()
	if err != nil {
		app.Logger.Error(err.Error())
		c.JSON(parseError(err))
		c.Abort()
	}
}

func parseError(err error) (int, gin.H) {
	switch err.(type) {
	case usecase_errors.IPermissionError:
		return http.StatusForbidden, gin.H{"error": err.Error()}
	case usecase_errors.IAlreadyExistsError:
		return http.StatusConflict, gin.H{"error": err.Error()}
	case usecase_errors.IUnauthorizedError:
		return http.StatusUnauthorized, gin.H{"error": err.Error()}
	case usecase_errors.IBadRequestError:
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	case usecase_errors.INotFoundError:
		return http.StatusNotFound, gin.H{"error": err.Error()}
	}

	return http.StatusServiceUnavailable, gin.H{
		"error": "Internal server error",
	}
}
