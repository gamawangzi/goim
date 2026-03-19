package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type RoomHandler struct {
	connMgr *BucketManager
	roomMgr *RoomManager
}

func NewRoomHandler(connMgr *BucketManager, roomMgr *RoomManager) *RoomHandler {
	return &RoomHandler{
		connMgr: connMgr,
		roomMgr: roomMgr,
	}
}

// RoomRequest 房间操作请求
type RoomRequest struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

// RoomPushRequest 房间推送请求
type RoomPushRequest struct {
	RoomID  string `json:"room_id"`
	Message string `json:"message"`
}

// RoomResponse 房间操作响应
type RoomResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (rh *RoomHandler) HandleJoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rh.sendResponse(w, 1, "只支持POST请求")
		return
	}

	var req RoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.sendResponse(w, 1, "请求格式错误")
		return
	}

	// 获取用户的 Channel
	channel := rh.connMgr.Get(req.UserID)
	if channel == nil {
		rh.sendResponse(w, 1, "用户不在线")
		return
	}

	// 获取或创建房间
	room := rh.roomMgr.GetOrCreate(req.RoomID)

	// 加入房间
	if err := channel.JoinRoom(room); err != nil {
		rh.sendResponse(w, 1, "加入房间失败")
		return
	}

	log.Printf("[RoomHandler] 用户 %s 加入房间 %s", req.UserID, req.RoomID)
	rh.sendResponse(w, 0, "加入房间成功")
}

func (rh *RoomHandler) HandleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rh.sendResponse(w, 1, "只支持POST请求")
		return
	}

	var req RoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.sendResponse(w, 1, "请求格式错误")
		return
	}

	// 获取用户的 Channel
	channel := rh.connMgr.Get(req.UserID)
	if channel == nil {
		rh.sendResponse(w, 1, "用户不在线")
		return
	}

	// 离开房间
	channel.LeaveRoom(req.RoomID)

	// 尝试删除空房间
	rh.roomMgr.Delete(req.RoomID)

	log.Printf("[RoomHandler] 用户 %s 离开房间 %s", req.UserID, req.RoomID)
	rh.sendResponse(w, 0, "离开房间成功")
}

func (rh *RoomHandler) HandlePushRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rh.sendResponse(w, 1, "只支持POST请求")
		return
	}

	var req RoomPushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rh.sendResponse(w, 1, "请求格式错误")
		return
	}

	// 构造房间消息
	roomMsg := RoomMessage{
		RoomID:  req.RoomID,
		Content: req.Message,
		From:    "system",
	}
	body, _ := json.Marshal(roomMsg)

	msg := Message{
		Ver:  1,
		Op:   OpRoomMessage,
		Body: string(body),
	}

	msgData, _ := json.Marshal(msg)

	// 推送到房间
	if err := rh.roomMgr.PushRoom(req.RoomID, string(msgData)); err != nil {
		log.Printf("[RoomHandler] 房间推送失败: %v", err)
		rh.sendResponse(w, 1, "推送失败")
		return
	}

	log.Printf("[RoomHandler] 房间 %s 推送成功: %s", req.RoomID, req.Message)
	rh.sendResponse(w, 0, "推送成功")
}

func (rh *RoomHandler) sendResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	resp := RoomResponse{Code: code, Msg: msg}
	json.NewEncoder(w).Encode(resp)
}