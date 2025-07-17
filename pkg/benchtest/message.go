package benchtest

import (
	"bufio"
	"github.com/google/uuid"
	"github.com/sanmuyan/xpkg/xnet"
	"google.golang.org/protobuf/proto"
	"net"
)

//go:generate protoc --go_out=../ bt_message.proto

const (
	NewMessage = iota + 20000
)

const (
	ReadBufferSize = 1024
)

func Unmarshal(data []byte) (*BtMessage, error) {
	msg := new(BtMessage)
	return msg, proto.Unmarshal(data, msg)
}

func Marshal(msg *BtMessage) ([]byte, error) {
	return proto.Marshal(msg)
}

func GenerateMessage(requestID string) *BtMessage {
	return &BtMessage{
		Ctl:       NewMessage,
		RequestID: requestID,
	}
}

func GenerateRequestID() string {
	return uuid.NewString()
}

func WriteTCP(msg *BtMessage, conn net.Conn) error {
	bp, err := Marshal(msg)
	if err != nil {
		return err
	}
	be, err := xnet.Encode(bp)
	if err != nil {
		return err
	}
	_, err = conn.Write(be)
	return err
}

func ReadTCP(reader *bufio.Reader) (*BtMessage, error) {
	be, err := xnet.Decode(reader)
	if err != nil {
		return nil, err
	}
	return Unmarshal(be)
}
