package main

import "encoding/json"

// 操作类型常量
const (
	OpHeartbeat      = 2 // 心跳
	OpHeartbeatReply = 3 // 心跳回复
	OpMessage        = 5 // 消息推送
	OpRoomMessage    = 6 // 房间消息
	OpAuth           = 7 // 认证请求
	OpAuthReply      = 8 // 认证回复
)

// Message 消息结构
type Message struct {
	Ver  int32  `json:"ver"`  // 协议版本
	Op   int32  `json:"op"`   // 操作类型
	Body string `json:"body"` // 消息体（JSON字符串）
}

// AuthRequest 认证请求
type AuthRequest struct {
	Token string `json:"token"` // 用户token（这里简化为用户ID）
}

// AuthReply 认证响应
type AuthReply struct {
	Code int    `json:"code"` // 0-成功，其他-失败
	Msg  string `json:"msg"`  // 消息
}

// PushMessage 推送消息
type PushMessage struct {
	Content string `json:"content"` // 消息内容
}

// RoomMessage 房间消息
type RoomMessage struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
	From    string `json:"from"` // 发送者
}


// ToJSON 将消息转为JSON字节
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// ParseMessage 解析JSON为消息
func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
