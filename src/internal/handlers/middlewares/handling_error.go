package handler_middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	api_errors "libs/src/internal/usecase/errors"
	"net/http"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	err := c.Errors.Last()
	if err != nil {
		c.JSON(parseError(err))
		c.Abort()
	}
}

func parseError(err error) (int, gin.H) {
	if errors.Is(err, api_errors.ErrPasswordsDontMatch) {
		return http.StatusBadRequest, gin.H{
			"error": "Passwords don't match",
		}
	}
	if errors.Is(err, api_errors.ErrInvalidToken) {
		return http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		}
	}
	if errors.Is(err, api_errors.ErrUserAlreadyExists) {
		return http.StatusConflict, gin.H{
			"error": "Account with username or email already exists",
		}
	}
	if errors.Is(err, api_errors.ErrInvalidCredentials) {
		return http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		}
	}
	if errors.Is(err, api_errors.ErrAlreadyLoggedIn) {
		return http.StatusConflict, gin.H{
			"error": "User is already logged in",
		}
	}
	if errors.Is(err, api_errors.ErrInvalidSession) {
		return http.StatusUnauthorized, gin.H{
			"error": "Invalid session",
		}
	}
	if errors.Is(err, api_errors.ErrProfileNotFound) {
		return http.StatusNotFound, gin.H{
			"error": "Profile not found",
		}
	}
	if errors.Is(err, api_errors.ErrNeedLoginForChangeProfile) {
		return http.StatusUnauthorized, gin.H{
			"error": "You need to login to change your profile",
		}
	}
	if errors.Is(err, api_errors.ErrInvalidBody) {
		return http.StatusBadRequest, gin.H{
			"error": "Invalid body",
		}
	}
	if errors.Is(err, api_errors.ErrInvalidData) {
		return http.StatusBadRequest, gin.H{
			"error": "Invalid data in body",
		}
	}
	if errors.Is(err, api_errors.ErrPasswordLight) {
		return http.StatusBadRequest, gin.H{
			"error": "the password must have at least 1 lower case, 1 upper case, a number and 1 special character, and be longer than 8 characters",
		}
	}
	if errors.Is(err, api_errors.ErrInvalidCode) {
		return http.StatusBadRequest, gin.H{
			"error": "Invalid code",
		}
	}
	if errors.Is(err, api_errors.ErrNotLoggedIn) {
		return http.StatusUnauthorized, gin.H{
			"error": "User is not logged in",
		}
	}
	if errors.Is(err, api_errors.ErrUserNotFound) {
		return http.StatusNotFound, gin.H{
			"error": "User not found",
		}
	}

	return http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
	}
}
