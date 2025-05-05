package handler_api

import (
	"github.com/gin-gonic/gin"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
	"strconv"
)

// @Summary Send message
// @Description Send a message to a chat
// @Tags Messages
// @Accept json
// @Produce json
// @Param ChatId path int true "Chat ID"
// @Param user body dto.SendMessageRequest true "Data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/{ChatId}/message/send [post]
func SendMessage(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	caller := c.MustGet("user").(dto.UserDTO)

	chatId := c.Param("chat_id")
	chatIdInt, err := strconv.Atoi(chatId)

	if err != nil {
		c.Error(usecase_errors.BadRequestError{Msg: "Invalid chat ID"})
		return
	}

	var messageRequest dto.SendMessageRequest
	if err := c.ShouldBindJSON(&messageRequest); err != nil {
		c.Error(usecase_errors.BadRequestError{Msg: err.Error()})
		return
	}

	messageService := services.NewMessageService(app)
	messagePreview, err := messageService.SendMessage(c.Request.Context(), caller, messageRequest, int64(chatIdInt))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, messagePreview)
}
