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

	// 创建WebSocket服务器
	wsServer := NewWSServer(connMgr)

	// 创建推送处理器
	pushHandler := NewPushHandler(connMgr)

	// 注册路由
	http.HandleFunc("/ws", wsServer.HandleWebSocket)
	http.HandleFunc("/push", pushHandler.HandlePush)

	// 启动服务
	log.Println("WebSocket 服务: ws://localhost:3102/ws")
	log.Println("HTTP 推送接口: http://localhost:3111/push")

	go func() {
		log.Fatal(http.ListenAndServe(":3102", nil))
	}()

	log.Fatal(http.ListenAndServe(":3111", nil))
}
