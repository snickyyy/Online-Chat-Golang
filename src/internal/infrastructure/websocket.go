package infrastructure

import (
	"github.com/gorilla/websocket"
	"libs/src/internal/dto"
)

type WebSocketClient struct {
	UserDto    *dto.UserDTO
	Conn       *websocket.Conn
	Messagebox chan []byte
}

type Message struct {
	From *WebSocketClient
	To   *WebSocketClient
	Data []byte
}

type WebSocketHub struct {
	Clients    map[*WebSocketClient]bool
	Messagebox chan []Message
	Add        chan *WebSocketClient
	Delete     chan *WebSocketClient
}
