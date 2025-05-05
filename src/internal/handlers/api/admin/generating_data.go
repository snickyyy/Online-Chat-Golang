package handler_api

import (
	"github.com/gin-gonic/gin"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase/admin"
	"libs/src/settings"
	"strconv"
)

// @Summary Generate users
// @Description generate users
// @Tags Admin
// @Accept json
// @Produce json
// @Param count query int false "count of users to generate"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/generate/user [post]
func GenerateUsers(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	count, _ := strconv.Atoi(c.Query("count"))
	if count == 0 {
		count = 50
	}

	generator := services.NewDataGenerator(app)
	err := generator.GenerateUsers(user, count)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, dto.MessageResponse{Message: "success"})
}

// @Summary Generate chats
// @Description generate chats
// @Tags Admin
// @Accept json
// @Produce json
// @Param count query int false "count of chats to generate"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/generate/chat [post]
func GenerateChats(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	count, _ := strconv.Atoi(c.Query("count"))
	if count == 0 {
		count = 50
	}

	generator := services.NewDataGenerator(app)
	err := generator.GenerateChats(user, count)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, dto.MessageResponse{Message: "success"})
}

// @Summary Generate members of chat
// @Description generate members of chat
// @Tags Admin
// @Accept json
// @Produce json
// @Param count query int false "count of members to generate"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/generate/members [post]
func GenerateChatMembers(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	count, _ := strconv.Atoi(c.Query("count"))
	if count == 0 {
		count = 50
	}

	generator := services.NewDataGenerator(app)
	err := generator.GenerateChatMembers(user, count)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, dto.MessageResponse{Message: "success"})
}

// @Summary Generate messages
// @Description generate messages
// @Tags Admin
// @Accept json
// @Produce json
// @Param count query int false "count of messages to generate"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/generate/message [post]
func GenerateMessages(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	count, _ := strconv.Atoi(c.Query("count"))
	if count == 0 {
		count = 50
	}

	generator := services.NewDataGenerator(app)
	err := generator.GenerateMessages(user, count)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, dto.MessageResponse{Message: "success"})
}
