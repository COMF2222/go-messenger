package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Connection struct {
	Conn *websocket.Conn
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &Connection{Conn: c}, nil
}

func (c *Connection) ReadMessage() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

func (c *Connection) WriteMessage(messageType int, data []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.Conn.WriteMessage(messageType, data)
}

func (c *Connection) Close() error {
	return c.Conn.Close()
}
