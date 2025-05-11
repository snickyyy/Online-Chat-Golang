package websocket

import (
	"github.com/gin-gonic/gin"
	"libs/src/internal/dto"
	"libs/src/internal/infrastructure"
	"libs/src/settings"
)

func Chat(hub *infrastructure.WebSocketHub, c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	conn, err := app.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	defer conn.Close()
	if err != nil {
		c.Error(err)
		return
	}

	client := infrastructure.NewWebSocketClient(conn, &user)

	hub.Add <- client
	go client.WritePump()

	client.ReadPump(hub)
}
