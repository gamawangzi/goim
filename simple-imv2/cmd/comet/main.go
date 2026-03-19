package main

import (
	"log"
	"net"
	"net/http"

	"simple-imv2/api/comet"
	cometServer "simple-imv2/internal/comet"

	"google.golang.org/grpc"
)

func main() {
	log.Println("=== Comet 服务启动 ===")

	bucketMgr := cometServer.NewBucketManager(32)
	wsServer := cometServer.NewWSServer(bucketMgr)
	grpcServer := cometServer.NewGrpcServer(bucketMgr)

	go func() {
		http.HandleFunc("/ws", wsServer.HandleWebSocket)
		log.Println("WebSocket 服务: ws://localhost:3102/ws")
		log.Fatal(http.ListenAndServe(":3102", nil))
	}()

	lis, err := net.Listen("tcp", ":3109")
	if err != nil {
		log.Fatalf("gRPC 监听失败: %v", err)
	}

	s := grpc.NewServer()
	comet.RegisterCometServiceServer(s, grpcServer)
	log.Println("gRPC 服务: localhost:3109")
	log.Fatal(s.Serve(lis))
}
