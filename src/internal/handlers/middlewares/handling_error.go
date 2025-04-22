package handler_middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"libs/src/internal/repositories"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
	"net/http"
)

func ErrorHandler(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	c.Next()
	err := c.Errors.Last()
	if err != nil {
		app.Logger.Error(fmt.Sprintf("\nrequest error: %v\nurl: %s\nmethod: %s\n", err.Error(), c.Request.URL.Path, c.Request.Method))
		c.JSON(parseError(err.Err))
		c.Abort()
	}
}

func parseError(err error) (int, gin.H) {
	if _, ok := err.(usecase_errors.IPermissionError); ok {
		return http.StatusForbidden, gin.H{"error": err.Error()}
	}
	if _, ok := err.(usecase_errors.IAlreadyExistsError); ok {
		return http.StatusConflict, gin.H{"error": err.Error()}
	}
	if _, ok := err.(usecase_errors.IUnauthorizedError); ok {
		return http.StatusUnauthorized, gin.H{"error": err.Error()}
	}
	if _, ok := err.(usecase_errors.IBadRequestError); ok {
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	}
	if _, ok := err.(usecase_errors.INotFoundError); ok {
		return http.StatusNotFound, gin.H{"error": err.Error()}
	}
	if errors.Is(err, repositories.ErrTimeout) {
		return http.StatusServiceUnavailable, gin.H{"error": "Request timeout"}
	}
	return http.StatusServiceUnavailable, gin.H{
		"error": "Internal server error",
	}
}
