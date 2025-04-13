package handler_api

import (
	"github.com/gin-gonic/gin"
	"libs/src/internal/domain/enums"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
	"net/http"
	"strconv"
)

// @Summary Create chat
// @Description Creating a new chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param user body dto.CreateChatRequest true "Data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/create [post]
func CreateChat(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	var chatData dto.CreateChatRequest

	if err := c.ShouldBindJSON(&chatData); err != nil {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewChatService(app)

	chat, err := service.CreateChat(chatData, user)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, chat)
}

// @Summary Invite to chat
// @Description inviting a user to an existing chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param invitee query string true "Invitee username"
// @Param chat_id query int true "chat id to invite to"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/invite [post]
func InviteToChat(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	inviter := c.MustGet("user").(dto.UserDTO)

	if inviter.Role == enums.ANONYMOUS || !inviter.IsActive {
		c.Error(api_errors.ErrUnauthorized)
		return
	}

	invitee := c.Query("invitee")
	chatID := c.Query("chat_id")
	if invitee == "" || chatID == "" {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewChatMemberService(app)

	chatIDInt, err := strconv.Atoi(chatID)
	if err != nil {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	err = service.InviteToChat(&inviter, invitee, int64(chatIDInt))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}

// @Summary Delete chat
// @Description deleting a chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param chatId path int true "Chat ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/delete/{chatId} [delete]
func DeleteChat(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	if user.Role == enums.ANONYMOUS || !user.IsActive {
		c.Error(api_errors.ErrUnauthorized)
		return
	}

	chatID := c.Param("chat_id")
	if chatID == "" {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewChatService(app)

	chatIDInt, err := strconv.Atoi(chatID)

	err = service.DeleteChat(user, int64(chatIDInt))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}
