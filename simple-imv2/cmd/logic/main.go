package main

import (
	"log"
	"net/http"

	"simple-imv2/internal/logic"
)

func main() {
	log.Println("=== Logic 服务启动 ===")

	cometClient, err := logic.NewCometClient("localhost:3109")
	if err != nil {
		log.Fatalf("连接 Comet 失败: %v", err)
	}
	defer cometClient.Close()

	roomMgr := logic.NewRoomManager()
	handler := logic.NewHandler(roomMgr, cometClient)

	http.HandleFunc("/room/join", handler.HandleJoinRoom)
	http.HandleFunc("/room/leave", handler.HandleLeaveRoom)
	http.HandleFunc("/room/push", handler.HandlePushRoom)

	log.Println("HTTP 服务: http://localhost:3111")
	log.Fatal(http.ListenAndServe(":3111", nil))
}
