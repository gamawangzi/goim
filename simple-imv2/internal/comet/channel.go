package comet

import (
	"github.com/gorilla/websocket"
)

type Channel struct {
	userID string
	conn   *websocket.Conn
}

func NewChannel(userID string, conn *websocket.Conn) *Channel {
	return &Channel{
		userID: userID,
		conn:   conn,
	}
}

func (c *Channel) Push(msg string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *Channel) Close() {
	c.conn.Close()
}
