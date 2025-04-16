package handler_api

import (
	"github.com/gin-gonic/gin"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
	"net/http"
	"strconv"
)

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
