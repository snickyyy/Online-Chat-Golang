package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"libs/src/internal/domain/enums"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	services "libs/src/internal/usecase"
	"libs/src/settings"
	"net/http"
)

func (suite *AppTestSuite) TestCreateChat() {
	url := "http://127.0.0.1:8000/messenger/chat/create"

	userService := services.NewUserService(settings.AppVar)
	authService := services.NewAuthService(settings.AppVar)

	err := userService.CreateSuperUser(suite.Ctx, "TestCreateChat", "TestCreateChat@test.com", "test123")
	suite.NoError(err)
	sess, err := authService.Login(suite.Ctx, dto.UserDTO{ID: 1, Role: enums.ANONYMOUS, IsActive: false}, dto.LoginRequest{UsernameOrEmail: "TestCreateChat", Password: "test123"})
	suite.NoError(err)

	dataCreateChat, _ := json.Marshal(dto.CreateChatRequest{
		Title:       "TestCreateChat",
		Description: "TestCreateChat",
	})
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(dataCreateChat))
	requestResult, err := suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, requestResult.StatusCode)

	request, _ = http.NewRequest("POST", url, bytes.NewBuffer(dataCreateChat))
	request.AddCookie(&http.Cookie{Name: "sessionID", Value: sess})
	requestResult, err = suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusOK, requestResult.StatusCode)
}

func (suite *AppTestSuite) TestInviteToChat() {
	chatCreateUrl := "http://127.0.0.1:8000/messenger/chat/create"
	inviteUrl := "http://127.0.0.1:8000/messenger/chat/invite?invitee=%s&chat_id=%d"

	// Setup services
	userService := services.NewUserService(settings.AppVar)
	chatService := services.NewChatService(settings.AppVar)
	authService := services.NewAuthService(settings.AppVar)

	// Create test users
	err := userService.CreateSuperUser(suite.Ctx, "TestInviter", "TestInviter@test.com", "test123")
	suite.NoError(err)
	err = userService.CreateSuperUser(suite.Ctx, "TestInvitee", "TestInvitee@test.com", "test123")
	suite.NoError(err)

	// Create session for the inviter
	sess, err := authService.Login(suite.Ctx, dto.UserDTO{ID: 1, Role: enums.ANONYMOUS, IsActive: false}, dto.LoginRequest{UsernameOrEmail: "TestInviter", Password: "test123"})
	suite.NoError(err)

	// Create a chat
	dataCreateChat, _ := json.Marshal(dto.CreateChatRequest{Title: "TestChat", Description: "TestChat"})

	request, _ := http.NewRequest("POST", chatCreateUrl, bytes.NewBuffer(dataCreateChat))
	request.AddCookie(&http.Cookie{Name: "sessionID", Value: sess})
	requestResult, err := suite.client.Do(request)
	suite.NoError(err)
	suite.Equal(http.StatusOK, requestResult.StatusCode)

	// Decode the response to get the chat ID
	var chatResponse dto.ChatDTO
	err = json.NewDecoder(requestResult.Body).Decode(&chatResponse)
	suite.NoError(err)
	chatID := chatResponse.ID

	// Invite the invitee to the chat
	inviteRequest, _ := http.NewRequest("POST", fmt.Sprintf(inviteUrl, "TestInvitee", chatID), nil)
	inviteRequest.AddCookie(&http.Cookie{Name: "sessionID", Value: sess})
	inviteResult, err := suite.client.Do(inviteRequest)
	suite.NoError(err)
	suite.Equal(http.StatusOK, inviteResult.StatusCode)

	// Check if the invitee is now a member of the chat
	userRepo := repositories.NewUserRepository(settings.AppVar)
	invitee, err := userRepo.GetByUsername(suite.Ctx, "TestInvitee")
	suite.NoError(err)

	inviteeInfo, err := chatService.ChatMemberRepository.GetMemberInfo(suite.Ctx, invitee.ID, chatID)
	suite.NoError(err)
	suite.Equal(inviteeInfo.MemberRole, dto.ChatMemberDTO{MemberRole: 0}.MemberRole)
}
