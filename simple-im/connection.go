package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*websocket.Conn // key: userID, value: websocket连接
	mu          sync.RWMutex               // 读写锁
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*websocket.Conn),
	}
}

// Add 添加连接
func (cm *ConnectionManager) Add(userID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[userID] = conn
	log.Printf("[ConnectionManager] 用户 %s 已连接，当前在线: %d", userID, len(cm.connections))
}

// Remove 移除连接
func (cm *ConnectionManager) Remove(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, userID)
	log.Printf("[ConnectionManager] 用户 %s 已断开，当前在线: %d", userID, len(cm.connections))
}

// Get 获取连接
func (cm *ConnectionManager) Get(userID string) (*websocket.Conn, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conn, ok := cm.connections[userID]
	return conn, ok
}

// Push 向指定用户推送消息
func (cm *ConnectionManager) Push(userID string, msg Message) error {
	conn, ok := cm.Get(userID)
	if !ok {
		return ErrUserNotOnline
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

// Count 获取在线用户数
func (cm *ConnectionManager) Count() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.connections)
}
