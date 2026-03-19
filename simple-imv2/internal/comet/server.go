package comet

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// WSServer WebSocket服务器
type WSServer struct {
	connMgr *BucketManager 
}

// NewWSServer 创建WebSocket服务器
func NewWSServer(connMgr *BucketManager ) *WSServer {
	return &WSServer{
		connMgr: connMgr,
	}
}

// HandleWebSocket 处理WebSocket连接
func (s *WSServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WSServer] 升级WebSocket失败: %v", err)
		return
	}

	log.Printf("[WSServer] 新连接来自: %s", r.RemoteAddr)

	// 等待认证
	userID, err := s.authenticate(conn)
	if err != nil {
		log.Printf("[WSServer] 认证失败: %v", err)
		conn.Close()
		return
	}

	// 认证成功，添加到连接管理器
	channel :=  NewChannel(userID,conn)
	s.connMgr.Add(userID, channel)
	defer s.connMgr.Remove(userID)
	defer conn.Close()

	// 读取消息循环
	s.readLoop(conn, userID)
}

// authenticate 处理认证
func (s *WSServer) authenticate(conn *websocket.Conn) (string, error) {
	// 设置认证超时
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer conn.SetReadDeadline(time.Time{}) // 清除超时

	_, data, err := conn.ReadMessage()
	if err != nil {
		return "", err
	}

	msg, err := ParseMessage(data)
	if err != nil {
		return "", err
	}

	// 检查是否是认证消息
	if msg.Op != OpAuth {
		s.sendAuthReply(conn, 1, "需要先认证")
		return "", ErrAuthFailed
	}

	// 解析认证请求
	var authReq AuthRequest
	if err := json.Unmarshal([]byte(msg.Body), &authReq); err != nil {
		s.sendAuthReply(conn, 1, "认证数据格式错误")
		return "", err
	}

	// 简单验证：token非空即通过
	if authReq.Token == "" {
		s.sendAuthReply(conn, 1, "token不能为空")
		return "", ErrAuthFailed
	}

	// 认证成功
	s.sendAuthReply(conn, 0, "认证成功")
	log.Printf("[WSServer] 用户 %s 认证成功", authReq.Token)
	return authReq.Token, nil
}

// sendAuthReply 发送认证回复
func (s *WSServer) sendAuthReply(conn *websocket.Conn, code int, msg string) {
	reply := AuthReply{Code: code, Msg: msg}
	body, _ := json.Marshal(reply)

	authMsg := Message{
		Ver:  1,
		Op:   OpAuthReply,
		Body: string(body),
	}

	data, _ := authMsg.ToJSON()
	conn.WriteMessage(websocket.TextMessage, data)
}

// readLoop 读取消息循环
func (s *WSServer) readLoop(conn *websocket.Conn, userID string) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WSServer] 用户 %s 读取消息失败: %v", userID, err)
			break
		}

		msg, err := ParseMessage(data)
		if err != nil {
			log.Printf("[WSServer] 解析消息失败: %v", err)
			continue
		}

		s.handleMessage(conn, userID, msg)
	}
}

// handleMessage 处理消息
func (s *WSServer) handleMessage(conn *websocket.Conn, userID string, msg *Message) {
	switch msg.Op {
	case OpHeartbeat:
		// 心跳响应
		s.sendHeartbeatReply(conn)
		log.Printf("[WSServer] 用户 %s 心跳", userID)
	default:
		log.Printf("[WSServer] 用户 %s 未知操作: %d", userID, msg.Op)
	}
}

// sendHeartbeatReply 发送心跳回复
func (s *WSServer) sendHeartbeatReply(conn *websocket.Conn) {
	msg := Message{Ver: 1, Op: OpHeartbeatReply, Body: ""}
	data, _ := msg.ToJSON()
	conn.WriteMessage(websocket.TextMessage, data)
}
