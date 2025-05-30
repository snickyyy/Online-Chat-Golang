package handler_api

import (
	"github.com/gin-gonic/gin"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	usecase_errors "libs/src/internal/usecase/errors"
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
		c.Error(usecase_errors.BadRequestError{Msg: err.Error()})
		return
	}

	service := services.NewChatService(app)

	chat, err := service.CreateChat(c.Request.Context(), chatData, user)
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
		c.Error(usecase_errors.BadRequestError{Msg: "chat_id is empty"})
		return
	}

	service := services.NewChatService(app)

	chatIDInt, _ := strconv.Atoi(chatID)

	err := service.DeleteChat(c.Request.Context(), user, int64(chatIDInt))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}

// @Summary Change chat
// @Description change chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param chatId path int true "Chat ID"
// @Param data body dto.ChangeChatRequest true "Data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/edit/{chatId} [patch]
func ChangeChat(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	chatID := c.Param("chat_id")
	if chatID == "" {
		c.Error(usecase_errors.BadRequestError{Msg: "chat_id is empty"})
		return
	}

	var chatData dto.ChangeChatRequest

	if err := c.ShouldBindJSON(&chatData); err != nil {
		c.Error(usecase_errors.BadRequestError{Msg: err.Error()})
		return
	}

	service := services.NewChatService(app)

	chatIDInt, _ := strconv.Atoi(chatID)

	err := service.ChangeChat(c.Request.Context(), user, int64(chatIDInt), chatData)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}

// @Summary Get chats for user
// @Description get all the chats in which the user consists
// @Tags Chat
// @Accept json
// @Produce json
// @Param search query string false "Search name"
// @Param page query int false "Page"
// @Success 200 {object} dto.ChatsForUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/all [get]
func GetChatsForUser(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	page := c.Query("page")
	search := c.Query("search")

	if page == "" {
		page = "1"
	}
	pageInt, _ := strconv.Atoi(page)

	service := services.NewChatService(app)

	var (
		chats []dto.ChatDTO
		err   error
	)

	if search != "" {
		chats, err = service.Search(c.Request.Context(), user, search, pageInt)
	} else {
		chats, err = service.GetListForUser(c.Request.Context(), user, pageInt)
	}

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, chats)
}

// @Summary Get chat info
// @Description get chat info by id
// @Tags Chat
// @Accept json
// @Produce json
// @Param ChatId path string true "Chat id"
// @Success 200 {object} dto.ChatDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/{ChatId} [get]
func GetChatInfo(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	chatID := c.Param("chat_id")
	if chatID == "" {
		c.Error(usecase_errors.BadRequestError{Msg: "chat_id is empty"})
		return
	}

	service := services.NewChatService(app)
	chatIDInt, _ := strconv.Atoi(chatID)
	chat, err := service.GetById(c.Request.Context(), user, int64(chatIDInt))

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chat)
}
