/*
 * @Author: wangqian
 * @Date: 2026-03-07 20:48:17
 * @LastEditors: wangqian
 * @LastEditTime: 2026-03-11 14:43:46
 */
package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Channel struct {
	userID string
	conn   *websocket.Conn
	rooms  map[string]*Room
	lock   sync.RWMutex
}

func NewChannel(userID string, conn *websocket.Conn) *Channel {
	return &Channel{
		userID: userID,
		conn:   conn,
		rooms:  make(map[string]*Room),
		lock:   sync.RWMutex{},
	}
}
func (c *Channel) Push(msg string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *Channel) JoinRoom(room *Room) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.rooms[room.id] = room
	room.Join(c.userID, c)
	return nil
}

func (c *Channel) LeaveRoom(roomID string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if room, ok := c.rooms[roomID]; ok {
		delete(c.rooms, roomID)
		room.Leave(c.userID)
	}
}

func (c *Channel) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for key, value := range c.rooms {
		delete(c.rooms, key)
		value.Leave(c.userID)
	}
}
