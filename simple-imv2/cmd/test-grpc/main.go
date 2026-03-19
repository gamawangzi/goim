package main

import (
	"context"
	"log"
	"time"

	pb "simple-imv2/api/comet"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:3109", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewCometServiceClient(conn)

	log.Println("=== 测试 PushMsg ===")
	pushResp, err := client.PushMsg(context.Background(), &pb.PushMsgReq{
		UserId:  "user123",
		Message: `{"ver":1,"op":5,"body":"{\"content\":\"Hello from gRPC\"}"}`,
	})
	if err != nil {
		log.Printf("PushMsg 错误: %v", err)
	} else {
		log.Printf("PushMsg 响应: code=%d, msg=%s", pushResp.Code, pushResp.Msg)
	}

	time.Sleep(1 * time.Second)

	log.Println("=== 测试 PushRoom ===")
	roomResp, err := client.PushRoom(context.Background(), &pb.PushRoomReq{
		RoomId:  "room001",
		Message: `{"ver":1,"op":6,"body":"{\"room_id\":\"room001\",\"content\":\"Room message\",\"from\":\"system\"}"}`,
		UserIds: []string{"user123", "user456"},
	})
	if err != nil {
		log.Printf("PushRoom 错误: %v", err)
	} else {
		log.Printf("PushRoom 响应: code=%d, msg=%s", roomResp.Code, roomResp.Msg)
	}
}
