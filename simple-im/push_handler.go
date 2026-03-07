/*
 * @Author: wangqian
 * @Date: 2026-03-03 20:01:24
 * @LastEditors: wangqian
 * @LastEditTime: 2026-03-05 14:32:40
 */
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// PushHandler HTTP推送处理器
type PushHandler struct {
	connMgr *BucketManager
}

// NewPushHandler 创建推送处理器
func NewPushHandler(connMgr *BucketManager) *PushHandler {
	return &PushHandler{connMgr: connMgr}
}

// PushRequest 推送请求
type PushRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// PushResponse 推送响应
type PushResponse struct {
	Code int    `json:"code"` // 0-成功，其他-失败
	Msg  string `json:"msg"`
}

// HandlePush 处理推送请求
func (h *PushHandler) HandlePush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendResponse(w, 1, "只支持POST请求")
		return
	}

	var req PushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, 1, "请求格式错误")
		return
	}

	// 构造推送消息
	pushMsg := PushMessage{Content: req.Message}
	body, _ := json.Marshal(pushMsg)

	msg := Message{
		Ver:  1,
		Op:   OpMessage,
		Body: string(body),
	}

	// 推送消息
	if err := h.connMgr.Push(req.UserID, msg); err != nil {
		log.Printf("[PushHandler] 推送失败: %v", err)
		h.sendResponse(w, 1, err.Error())
		return
	}

	log.Printf("[PushHandler] 推送成功: %s -> %s", req.UserID, req.Message)
	h.sendResponse(w, 0, "推送成功")
}

// sendResponse 发送响应
func (h *PushHandler) sendResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	resp := PushResponse{Code: code, Msg: msg}
	json.NewEncoder(w).Encode(resp)
}
