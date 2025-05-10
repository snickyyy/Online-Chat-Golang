package infrastructure

import (
	"libs/src/settings"
	"sync"
)

type Message struct {
	From   *WebSocketClient
	ToChat int64
	Data   string
}

type WebSocketHub struct {
	App         *settings.App
	Mx          *sync.RWMutex
	Clients     map[*WebSocketClient]bool //todo хранить соединения в редисе
	Messagebox  chan *Message
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
		Messagebox:  make(chan *Message),
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
