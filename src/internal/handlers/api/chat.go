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

	chatID := c.Param("chat_id")
	if chatID == "" {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewChatService(app)

	chatIDInt, _ := strconv.Atoi(chatID)

	err := service.DeleteChat(user, int64(chatIDInt))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}
