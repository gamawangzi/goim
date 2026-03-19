/*
 * @Author: wangqian
 * @Date: 2026-03-03 20:01:41
 * @LastEditors: wangqian
 * @LastEditTime: 2026-03-05 14:31:32
 */
package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("=== Simple IM Server 启动 ===")

	// 创建连接管理器
	connMgr := NewBucketManager(32)

	// 创建房间管理器
	roomMgr := NewRoomManager()

	// 创建WebSocket服务器
	wsServer := NewWSServer(connMgr)

	// 创建推送处理器
	pushHandler := NewPushHandler(connMgr)

	// 创建房间处理器
	roomHandler := NewRoomHandler(connMgr, roomMgr)

	// 注册路由
	http.HandleFunc("/ws", wsServer.HandleWebSocket)
	http.HandleFunc("/push", pushHandler.HandlePush)
	http.HandleFunc("/room/join", roomHandler.HandleJoinRoom)
	http.HandleFunc("/room/leave", roomHandler.HandleLeaveRoom)
	http.HandleFunc("/room/push", roomHandler.HandlePushRoom)

	// 启动服务
	log.Println("WebSocket 服务: ws://localhost:3102/ws")
	log.Println("HTTP 推送接口: http://localhost:3111/push")
	log.Println("房间接口: http://localhost:3111/room/*")

	go func() {
		log.Fatal(http.ListenAndServe(":3102", nil))
	}()

	log.Fatal(http.ListenAndServe(":3111", nil))
}
