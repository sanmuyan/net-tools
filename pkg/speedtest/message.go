package speedtest

import (
	"bufio"
	"encoding/binary"
	"github.com/quic-go/quic-go"
	"github.com/sanmuyan/xpkg/xnet"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
)

//go:generate protoc --go_out=../ st_message.proto

const (
	NewTest = iota + 10000
	NewData
)

const (
	ReadBufferSize = 1024 * 1024
	TCPDataSize    = 1024 * 128
	QUICDataSize   = 1024 * 128
)

func Unmarshal(data []byte) (*StMessage, error) {
	msg := new(StMessage)
	return msg, proto.Unmarshal(data, msg)
}

func Marshal(msg *StMessage) ([]byte, error) {
	return proto.Marshal(msg)
}

func WriteTCP(msg *StMessage, conn net.Conn) error {
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

func ReadTCP(reader *bufio.Reader) (*StMessage, error) {
	be, err := xnet.Decode(reader)
	if err != nil {
		return nil, err
	}
	return Unmarshal(be)
}

func ReadQUIC(stream *quic.Stream) (*StMessage, error) {
	var msgLen uint32
	err := binary.Read(stream, binary.BigEndian, &msgLen)
	if err != nil {
		return nil, err
	}
	data := make([]byte, msgLen)
	_, err = io.ReadFull(stream, data)
	if err != nil {
		return nil, err
	}
	return Unmarshal(data)
}

func WriteQUIC(msg *StMessage, stream *quic.Stream) error {
	bp, err := Marshal(msg)
	if err != nil {
		return err
	}
	_ = binary.Write(stream, binary.BigEndian, uint32(len(bp)))
	_, err = stream.Write(bp)
	return err
}

var PreMessageTCP = make([]byte, TCPDataSize)

func init() {
	for i := range PreMessageTCP {
		PreMessageTCP[i] = 'x'
	}
}

var PreMessageQUIC = make([]byte, QUICDataSize)

func init() {
	for i := range PreMessageQUIC {
		PreMessageQUIC[i] = 'x'
	}
}
