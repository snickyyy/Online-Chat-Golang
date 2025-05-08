package infrastructure

import (
	"fmt"
	"github.com/gorilla/websocket"
	"libs/src/internal/dto"
	"libs/src/settings"
	"sync"
)

type WebSocketClient struct {
	Mx            *sync.RWMutex
	UserDto       *dto.UserDTO
	Subscriptions map[int64]bool
	Conn          *websocket.Conn
	Messagebox    chan []byte
}

func NewWebSocketClient(conn *websocket.Conn, user *dto.UserDTO) *WebSocketClient {
	return &WebSocketClient{
		Mx:            &sync.RWMutex{},
		UserDto:       user,
		Conn:          conn,
		Messagebox:    make(chan []byte),
		Subscriptions: make(map[int64]bool),
	}
}

func (c *WebSocketClient) Send(message []byte) {
	if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
		settings.AppVar.Logger.Error(fmt.Sprintf("websocket send message error %v", err))
	}
}

type Message struct {
	From   *WebSocketClient
	ToChat int64
	Data   []byte
}

type WebSocketHub struct {
	App         *settings.App
	Mx          *sync.RWMutex
	Clients     map[*WebSocketClient]bool //todo хранить соединения в редисе
	Messagebox  chan []*Message
	Add         chan *WebSocketClient
	Delete      chan *WebSocketClient
	Subscribe   chan *Subscribe
	Unsubscribe chan *Subscribe
}

func NewWebSocketHub(app *settings.App) *WebSocketHub {
	return &WebSocketHub{
		App:         app,
		Mx:          &sync.RWMutex{},
		Clients:     make(map[*WebSocketClient]bool),
		Messagebox:  make(chan []*Message),
		Add:         make(chan *WebSocketClient),
		Delete:      make(chan *WebSocketClient),
		Subscribe:   make(chan *Subscribe),
		Unsubscribe: make(chan *Subscribe),
	}
}

type Subscribe struct {
	Client *WebSocketClient
	ChatId int64
}
