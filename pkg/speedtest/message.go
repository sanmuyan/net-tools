package speedtest

import (
	"bufio"
	"github.com/sanmuyan/xpkg/xnet"
	"google.golang.org/protobuf/proto"
	"net"
)

//go:generate protoc --go_out=../ message.proto

const (
	NewTest = iota + 10000
	NewData
)

const (
	ReadBufferSize = 1024 * 1024
	TCPDataSize    = 1024 * 128
)

func Unmarshal(data []byte) (*Message, error) {
	msg := new(Message)
	return msg, proto.Unmarshal(data, msg)
}

func Marshal(msg *Message) ([]byte, error) {
	return proto.Marshal(msg)
}

func WriteTCP(msg *Message, conn net.Conn) error {
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

func ReadTCP(reader *bufio.Reader) (*Message, error) {
	be, err := xnet.Decode(reader)
	if err != nil {
		return nil, err
	}
	return Unmarshal(be)
}

var PreMessageTCP = make([]byte, TCPDataSize)

func init() {
	for i := range PreMessageTCP {
		PreMessageTCP[i] = 'x'
	}
}
