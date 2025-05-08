package websocket

import (
	"encoding/json"
	"go.uber.org/zap"
	"libs/src/internal/infrastructure"
	services "libs/src/internal/usecase"
)

func RunProcessHub(hub *infrastructure.WebSocketHub) {
	msgService := services.NewMessageService(hub.App)
	chatMembersService := services.NewChatMemberService(hub.App)

	for {
		select {

		case client := <-hub.Add:
			hub.Mx.Lock()
			hub.Clients[client] = true
			hub.Mx.Unlock()

		case client := <-hub.Delete:
			hub.Mx.Lock()
			delete(hub.Clients, client)
			hub.Mx.Unlock()

		case subscribe := <-hub.Subscribe:
			if chatMembersService.IsInChat(hub.App.Ctx, subscribe.Client.UserDto.ID, subscribe.ChatId) {
				subscribe.Client.Mx.Lock()
				subscribe.Client.Subscriptions[subscribe.ChatId] = true
				subscribe.Client.Mx.Unlock()
			}

		case unsubscribe := <-hub.Unsubscribe:
			unsubscribe.Client.Mx.Lock()
			delete(unsubscribe.Client.Subscriptions, unsubscribe.ChatId)
			unsubscribe.Client.Mx.Unlock()

		case messages := <-hub.Messagebox:
			for _, msg := range messages {
				messagePreview, err := msgService.NewMessage(hub.App.Ctx, msg.From.UserDto, string(msg.Data), msg.ToChat)
				if err != nil {
					hub.App.Logger.Error("new message failed", zap.Error(err))
					continue
				}

				messagePreviewToJson, _ := json.Marshal(messagePreview)

				for k, _ := range hub.Clients {
					if k.Subscriptions[msg.ToChat] {
						k.Messagebox <- messagePreviewToJson
					}
				}
			}
		}
	}
}
