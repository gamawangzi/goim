package comet

import (
	"context"
	"log"

	pb "simple-imv2/api/comet"
)

type GrpcServer struct {
	pb.UnimplementedCometServiceServer
	bucketMgr *BucketManager
}

func NewGrpcServer(bucketMgr *BucketManager) *GrpcServer {
	return &GrpcServer{
		bucketMgr: bucketMgr,
	}
}

func (s *GrpcServer) PushMsg(ctx context.Context, req *pb.PushMsgReq) (*pb.PushMsgReply, error) {
	log.Printf("[Comet gRPC] 推送消息给用户 %s", req.UserId)

	channel := s.bucketMgr.Get(req.UserId)
	if channel == nil {
		return &pb.PushMsgReply{Code: 1, Msg: "用户不在线"}, nil
	}

	err := channel.Push(req.Message)
	if err != nil {
		return &pb.PushMsgReply{Code: 1, Msg: "推送失败"}, nil
	}

	return &pb.PushMsgReply{Code: 0, Msg: "推送成功"}, nil
}

func (s *GrpcServer) PushRoom(ctx context.Context, req *pb.PushRoomReq) (*pb.PushRoomReply, error) {
	log.Printf("[Comet gRPC] 推送房间消息到 %s, 用户数: %d", req.RoomId, len(req.UserIds))

	successCount := 0
	for _, userID := range req.UserIds {
		channel := s.bucketMgr.Get(userID)
		if channel != nil {
			err := channel.Push(req.Message)
			if err == nil {
				successCount++
			}
		}
	}

	log.Printf("[Comet gRPC] 房间 %s 推送完成: %d/%d", req.RoomId, successCount, len(req.UserIds))
	return &pb.PushRoomReply{Code: 0, Msg: "推送完成"}, nil
}