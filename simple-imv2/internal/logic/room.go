package logic

import (
	"sync"
)

type Room struct {
	id      string
	userIDs map[string]bool
	lock    sync.RWMutex
}

type RoomManager struct {
	rooms map[string]*Room
	lock  sync.RWMutex
}

func NewRoom(roomID string) *Room {
	return &Room{
		id:      roomID,
		userIDs: make(map[string]bool),
	}
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) GetOrCreate(roomID string) *Room {
	rm.lock.RLock()
	if room, ok := rm.rooms[roomID]; ok {
		rm.lock.RUnlock()
		return room
	}
	rm.lock.RUnlock()

	rm.lock.Lock()
	defer rm.lock.Unlock()
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
	if room, ok := rm.rooms[roomID]; ok {
		if room.Count() == 0 {
			delete(rm.rooms, roomID)
		}
	}
}

func (r *Room) Join(userID string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.userIDs[userID] = true
}

func (r *Room) Leave(userID string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.userIDs, userID)
}

func (r *Room) GetUserIDs() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	userIDs := make([]string, 0, len(r.userIDs))
	for userID := range r.userIDs {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

func (r *Room) Count() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.userIDs)
}
