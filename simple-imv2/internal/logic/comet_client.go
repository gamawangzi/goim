package logic

import (
	"context"
	"log"

	pb "simple-imv2/api/comet"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CometClient struct {
	client pb.CometServiceClient
	conn   *grpc.ClientConn
}

func NewCometClient(addr string) (*CometClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &CometClient{
		client: pb.NewCometServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *CometClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *CometClient) PushMsg(userID, message string) error {
	resp, err := c.client.PushMsg(context.Background(), &pb.PushMsgReq{
		UserId:  userID,
		Message: message,
	})
	if err != nil {
		return err
	}
	log.Printf("[Logic] PushMsg 响应: %s", resp.Msg)
	return nil
}

func (c *CometClient) PushRoom(roomID, message string, userIDs []string) error {
	resp, err := c.client.PushRoom(context.Background(), &pb.PushRoomReq{
		RoomId:  roomID,
		Message: message,
		UserIds: userIDs,
	})
	if err != nil {
		return err
	}
	log.Printf("[Logic] PushRoom 响应: %s", resp.Msg)
	return nil
}
