package comet

import "encoding/json"

const (
	OpHeartbeat      = 2
	OpHeartbeatReply = 3
	OpMessage        = 5
	OpRoomMessage    = 6
	OpAuth           = 7
	OpAuthReply      = 8
)

type Message struct {
	Ver  int32  `json:"ver"`
	Op   int32  `json:"op"`
	Body string `json:"body"`
}

type AuthRequest struct {
	Token string `json:"token"`
}

type AuthReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
