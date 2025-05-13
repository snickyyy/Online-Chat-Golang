package websocket

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"libs/src/internal/dto"
	"libs/src/internal/infrastructure"
	services "libs/src/internal/usecase"
	"libs/src/settings"
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
			client.Mx.Lock()
			close(client.Messagebox)
			client.Conn.Close()
			client.Mx.Unlock()
			hub.Mx.Lock()
			delete(hub.Clients, client)
			hub.Mx.Unlock()

		case subscribe := <-hub.Subscribe:
			if chatMembersService.IsInChat(hub.App.Ctx, subscribe.Client.UserDto.ID, subscribe.ChatId) {

				subscribe.Client.Mx.Lock()
				subscribe.Client.Subscriptions[subscribe.ChatId] = true
				subscribe.Client.Mx.Unlock()

				response, _ := json.Marshal(dto.MessageInfo{Message: "successfully subscribed"})
				subscribe.Client.Messagebox <- &dto.ChatCommunication{
					Action: hub.App.Config.WsConfig.Actions.Subscribe,
					Body:   response,
				}
			} else {
				sendErrorMessage(subscribe.Client, errors.New("you are not a member of this chat"))
			}

		case unsubscribe := <-hub.Unsubscribe:
			unsubscribe.Client.Mx.Lock()
			delete(unsubscribe.Client.Subscriptions, unsubscribe.ChatId)
			unsubscribe.Client.Mx.Unlock()

			response, _ := json.Marshal(dto.MessageInfo{Message: "successfully unsubscribed"})
			unsubscribe.Client.Messagebox <- &dto.ChatCommunication{
				Action: hub.App.Config.WsConfig.Actions.Unsubscribe,
				Body:   response,
			}

		case msg := <-hub.Messagebox:
			if !isSubscribed(msg.From, msg.ToChat) {
				sendErrorMessage(msg.From, errors.New("you are not subscribed to this chat"))
			}

			messagePreview, err := msgService.NewMessage(hub.App.Ctx, msg.From.UserDto, msg.Data, msg.ToChat)

			if err != nil {
				hub.App.Logger.Error("new message failed", zap.Error(err))
				sendErrorMessage(msg.From, err)
				continue
			}

			messageBody, _ := json.Marshal(
				dto.MessagePreview{
					ChatId:      msg.ToChat,
					MessageBody: messagePreview,
				},
			)
			preparedMsg := dto.ChatCommunication{
				Action: hub.App.Config.WsConfig.Actions.SendMessage,
				Body:   messageBody,
			}
			go sendMessageToClients(hub.Clients, msg.ToChat, &preparedMsg)
		}
	}
}

func sendErrorMessage(client *infrastructure.WebSocketClient, err error) {
	client.Mx.Lock()
	defer client.Mx.Unlock()
	errorMessageBody, _ := json.Marshal(dto.ErrorMessage{Error: err.Error()})
	preparedMsg := dto.ChatCommunication{Action: settings.AppVar.Config.WsConfig.Actions.ErrorMessage, Body: errorMessageBody}
	client.Messagebox <- &preparedMsg
}

func isSubscribed(client *infrastructure.WebSocketClient, chatId int64) bool {
	client.Mx.RLock()
	defer client.Mx.RUnlock()
	_, isSubscribed := client.Subscriptions[chatId]
	return isSubscribed
}

func sendMessageToClients(clients map[*infrastructure.WebSocketClient]bool, withChat int64, msg *dto.ChatCommunication) {
	for client := range clients {
		client.Mx.RLock()
		_, isSubscribed := client.Subscriptions[withChat]
		client.Mx.RUnlock()
		if isSubscribed {
			client.Mx.Lock()
			client.Messagebox <- msg
			client.Mx.Unlock()
		}
	}
}
