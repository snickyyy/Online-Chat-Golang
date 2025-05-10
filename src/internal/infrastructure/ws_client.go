package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"libs/src/internal/dto"
	"libs/src/settings"
	"sync"
)

func NewWebSocketClient(conn *websocket.Conn, user *dto.UserDTO) *WebSocketClient {
	return &WebSocketClient{
		Mx:            &sync.RWMutex{},
		UserDto:       user,
		Conn:          conn,
		Messagebox:    make(chan *dto.ChatCommunication),
		Subscriptions: make(map[int64]bool),
	}
}

type WebSocketClient struct {
	Mx            *sync.RWMutex
	UserDto       *dto.UserDTO
	Subscriptions map[int64]bool
	Conn          *websocket.Conn
	Messagebox    chan *dto.ChatCommunication
}

func (c *WebSocketClient) ReadPump(hub *WebSocketHub) {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			hub.App.Logger.Error(fmt.Sprintf("websocket read message error %v", err))
			break
		}
		if len(message) == 0 {
			continue
		}

		var errorMsg dto.ChatCommunication
		handleErr := func(err error) {
			hub.App.Logger.Error(fmt.Sprintf("websocket pasrse message error %v", err))
			errorBody, _ := json.Marshal(dto.ErrorMessage{Error: err.Error()})
			errorMsg = dto.ChatCommunication{Action: "error", Body: errorBody}
			c.Messagebox <- &errorMsg
		}

		parsedMsg := dto.ChatCommunication{}
		err = json.Unmarshal(message, &parsedMsg)

		if err != nil {
			handleErr(err)
			continue
		}

		switch parsedMsg.Action {
		case hub.App.Config.WsConfig.Actions.Subscribe:
			subscribe := dto.Subscribe{}
			err = json.Unmarshal(parsedMsg.Body, &subscribe)
			if err != nil {
				handleErr(err)
				continue
			}

			hub.Subscribe <- &Subscribe{Client: c, ChatId: subscribe.ChatId}
		case hub.App.Config.WsConfig.Actions.Unsubscribe:
			unsubscribe := dto.Subscribe{}
			err = json.Unmarshal(parsedMsg.Body, &unsubscribe)
			if err != nil {
				handleErr(err)
				continue
			}
			hub.Unsubscribe <- &Subscribe{Client: c, ChatId: unsubscribe.ChatId}
		case hub.App.Config.WsConfig.Actions.SendMessage:
			sendMessage := dto.SendMessageRequest{}
			err = json.Unmarshal(parsedMsg.Body, &sendMessage)
			if err != nil {
				handleErr(err)
				continue
			}
			hub.Messagebox <- &Message{From: c, ToChat: sendMessage.ChatId, Data: sendMessage.Message}
		default:
			handleErr(fmt.Errorf("unknown action %s", parsedMsg.Action))
		}
	}
}

func (c *WebSocketClient) send(message dto.ChatCommunication) {
	messageBytes, _ := json.Marshal(message)

	c.Mx.Lock()
	if err := c.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		settings.AppVar.Logger.Error(fmt.Sprintf("websocket send message error %v", err))
	}
	defer c.Mx.Unlock()
}
