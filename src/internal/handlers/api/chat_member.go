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
// @Tags ChatMembers
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

// @Summary Change member role
// @Description сhanges the role of the user if possible
// @Tags ChatMembers
// @Accept json
// @Produce json
// @Param chat_id path int true "chat in which you want to change the role"
// @Param member_username path string true "target member username"
// @Param NewRole body dto.ChangeMemberRoleRequest true "new role for the member"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/{chat_id}/members/{member_username}/change-role [patch]
func ChangeMemberRole(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	caller := c.MustGet("user").(dto.UserDTO)

	var newRole dto.ChangeMemberRoleRequest
	if err := c.ShouldBindJSON(&newRole); err != nil {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	chatId := c.Param("chat_id")
	memberUsername := c.Param("member_username")

	if chatId == "" || memberUsername == "" {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewChatMemberService(app)

	chatIdInt, _ := strconv.Atoi(chatId)

	err := service.ChangeMemberRole(caller, int64(chatIdInt), memberUsername, newRole.NewRole)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}

// @Summary Delete member
// @Description Remove the member from the chat
// @Tags ChatMembers
// @Accept json
// @Produce json
// @Param chat_id path int true "chat from which you want to delete a member"
// @Param member_username path string true "target member username"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/{chat_id}/members/{member_username}/delete [delete]
func DeleteMember(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	caller := c.MustGet("user").(dto.UserDTO)

	chatId, _ := strconv.Atoi(c.Param("chat_id"))
	memberUsername := c.Param("member_username")
	if chatId == 0 || memberUsername == "" {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewChatMemberService(app)

	err := service.DeleteMember(caller, int64(chatId), memberUsername)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}

// @Summary Get member list
// @Description Get member list of chat
// @Tags ChatMembers
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param chat_id path int true "Chat id to get members from"
// @Success 200 {object} dto.MemberListPreview
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /messenger/chat/{chat_id}/members/all [get]
func GetMemberList(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	caller := c.MustGet("user").(dto.UserDTO)

	chatId, _ := strconv.Atoi(c.Param("chat_id"))
	page := c.Query("page")
	if page == "" {
		page = "1"
	}

	pageInt, _ := strconv.Atoi(page)

	if chatId == 0 {
		c.Error(api_errors.ErrInvalidData)
		return
	}

	service := services.NewChatMemberService(app)

	members, err := service.GetList(caller, int64(chatId), pageInt)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, members)
}
