package handler_middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	api_errors "libs/src/internal/usecase/errors"
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
	if errors.Is(err, api_errors.ErrInvalidPassword) {
		return http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		}
	}
	if errors.Is(err, api_errors.ErrUnauthorized) {
		return http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		}
	}
	if errors.Is(err, api_errors.ErrSamePassword) {
		return http.StatusBadRequest, gin.H{
			"error": "Passwords must be different",
		}
	}
	if errors.Is(err, api_errors.ErrChatAlreadyExists) {
		return http.StatusConflict, gin.H{
			"error": "Chat already exists",
		}
	}
	if errors.Is(err, api_errors.ErrChatNotFound) {
		return http.StatusNotFound, gin.H{
			"error": "Chat not found",
		}
	}
	if errors.Is(err, api_errors.ErrNotEnoughPermissionsForInviting) {
		return http.StatusForbidden, gin.H{
			"error": "not enough permissions for inviting",
		}
	}
	if errors.Is(err, api_errors.ErrUserAlreadyInChat) {
		return http.StatusBadRequest, gin.H{
			"error": "user already in chat or chat not exists",
		}
	}
	if errors.Is(err, api_errors.ErrUserNotInChat) {
		return http.StatusNotFound, gin.H{
			"error": "User not in chat",
		}
	}
	if errors.Is(err, api_errors.ErrNotEnoughPermissionsForChangeRole) {
		return http.StatusNotFound, gin.H{
			"error": "doesn't have enough permissions to change role",
		}
	}
	if errors.Is(err, api_errors.ErrInviterNotInChat) {
		return http.StatusBadRequest, gin.H{
			"error": "Inviter not in chat",
		}
	}
	if errors.Is(err, api_errors.ErrNotEnoughPermissionsForDelete) {
		return http.StatusForbidden, gin.H{
			"error": "not enough permissions for delete",
		}
	}
	if errors.Is(err, api_errors.ErrNotEnoughPermissionsForChangeChat) {
		return http.StatusForbidden, gin.H{
			"error": "not enough permissions for change chat",
		}
	}
	if errors.Is(err, api_errors.ErrNotEnoughPermissions) {
		return http.StatusForbidden, gin.H{
			"error": "not enough permissions",
		}
	}

	return http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
	}
}
