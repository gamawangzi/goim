/*
 * @Author: wangqian
 * @Date: 2026-03-07 20:40:23
 * @LastEditors: wangqian
 * @LastEditTime: 2026-03-11 14:42:27
 */
package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	id       string
	channels map[string]*Channel
	lock     sync.RWMutex
}

// 管理房间结构体
type RoomManager struct {
	rooms map[string]*Room
	lock  sync.RWMutex
}

func NewRoom(roomID string) *Room {
	return &Room{
		id:       roomID,
		channels: make(map[string]*Channel),
		lock:     sync.RWMutex{},
	}
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
		lock:  sync.RWMutex{},
	}
}

func (rm *RoomManager) Get(roomID string) *Room {
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	if val, ok := rm.rooms[roomID]; ok {
		return val
	}
	return nil
}

func (rm *RoomManager) GetOrCreate(roomID string) *Room {
	// 先尝试读取
	rm.lock.RLock()
	if room, ok := rm.rooms[roomID]; ok {
		rm.lock.RUnlock()
		return room
	}
	rm.lock.RUnlock()

	// 不存在则创建
	rm.lock.Lock()
	defer rm.lock.Unlock()
	// 双重检查
	if room, ok := rm.rooms[roomID]; ok {
		return room
	}
	room := NewRoom(roomID)
	rm.rooms[roomID] = room
	return room
}

func (rm *RoomManager) Delete(roomID string) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if val, ok := rm.rooms[roomID]; ok {
		if val.Count() == 0 {
			// 房间为空
			delete(rm.rooms, roomID)
		}
	}
}

func (rm *RoomManager) PushRoom(roomID string, msg string) error {
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	if val, ok := rm.rooms[roomID]; ok {
		return val.Push(msg)
	}
	return nil
}
func (r *Room) Join(userID string, channel *Channel) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.channels[userID] = channel
}

func (r *Room) Leave(userID string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.channels, userID)
}

func (r *Room) Push(msg string) error {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, channel := range r.channels {
		err := channel.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			continue // 忽略推送失败的用户
		}
	}
	return nil
}

func (r *Room) Count() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.channels)
}
