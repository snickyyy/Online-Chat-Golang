package websocket

import (
	"github.com/gin-gonic/gin"
	"libs/src/internal/domain/enums"
	"libs/src/internal/dto"
	"libs/src/internal/infrastructure"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
)

func Chat(hub *infrastructure.WebSocketHub, c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)
	if user.Role <= enums.ANONYMOUS || !user.IsActive {
		c.Error(usecase_errors.UnauthorizedError{Msg: "User is not authorized"})
		return

	}

	conn, err := app.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	defer conn.Close()
	if err != nil {
		c.Error(err)
		return
	}

	client := infrastructure.NewWebSocketClient(conn, &user)

	hub.Add <- client
	go client.WritePump(hub)

	client.ReadPump(hub)

	hub.Delete <- client
}
