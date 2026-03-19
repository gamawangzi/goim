package logic

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	roomMgr     *RoomManager
	cometClient *CometClient
}

func NewHandler(roomMgr *RoomManager, cometClient *CometClient) *Handler {
	return &Handler{
		roomMgr:     roomMgr,
		cometClient: cometClient,
	}
}

type JoinRoomReq struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

type PushRoomReq struct {
	RoomID  string `json:"room_id"`
	Message string `json:"message"`
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (h *Handler) HandleJoinRoom(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, 1, "请求格式错误")
		return
	}

	room := h.roomMgr.GetOrCreate(req.RoomID)
	room.Join(req.UserID)

	log.Printf("[Logic] 用户 %s 加入房间 %s", req.UserID, req.RoomID)
	h.sendResponse(w, 0, "加入成功")
}

func (h *Handler) HandleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, 1, "请求格式错误")
		return
	}

	room := h.roomMgr.GetOrCreate(req.RoomID)
	room.Leave(req.UserID)
	h.roomMgr.Delete(req.RoomID)

	log.Printf("[Logic] 用户 %s 离开房间 %s", req.UserID, req.RoomID)
	h.sendResponse(w, 0, "离开成功")
}

func (h *Handler) HandlePushRoom(w http.ResponseWriter, r *http.Request) {
	var req PushRoomReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendResponse(w, 1, "请求格式错误")
		return
	}

	room := h.roomMgr.GetOrCreate(req.RoomID)
	userIDs := room.GetUserIDs()

	if len(userIDs) == 0 {
		h.sendResponse(w, 0, "房间无用户")
		return
	}

	err := h.cometClient.PushRoom(req.RoomID, req.Message, userIDs)
	if err != nil {
		log.Printf("[Logic] 推送失败: %v", err)
		h.sendResponse(w, 1, "推送失败")
		return
	}

	log.Printf("[Logic] 房间 %s 推送成功", req.RoomID)
	h.sendResponse(w, 0, "推送成功")
}

func (h *Handler) sendResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Code: code, Msg: msg})
}
