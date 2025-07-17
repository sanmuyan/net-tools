package netbench

import (
	"bufio"
	"github.com/google/uuid"
	"github.com/sanmuyan/xpkg/xnet"
	"google.golang.org/protobuf/proto"
	"net"
)

//go:generate protoc --go_out=../ nb_message.proto

const (
	NewMessage = iota + 20000
)

const (
	ReadBufferSize = 1024
)

func Unmarshal(data []byte) (*NbMessage, error) {
	msg := new(NbMessage)
	return msg, proto.Unmarshal(data, msg)
}

func Marshal(msg *NbMessage) ([]byte, error) {
	return proto.Marshal(msg)
}

func GenerateMessage(requestID string) *NbMessage {
	return &NbMessage{
		Ctl:       NewMessage,
		RequestID: requestID,
	}
}

func GenerateRequestID() string {
	return uuid.NewString()
}

func WriteTCP(msg *NbMessage, conn net.Conn) error {
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

func ReadTCP(reader *bufio.Reader) (*NbMessage, error) {
	be, err := xnet.Decode(reader)
	if err != nil {
		return nil, err
	}
	return Unmarshal(be)
}
