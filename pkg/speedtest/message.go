package speedtest

import (
	"bufio"
	"google.golang.org/protobuf/proto"
	"net"
)

//go:generate protoc --go_out=../ message.proto

type Call func(msg *Message, err error) (exit bool)

type Message struct {
	*Options
}

func NewMessage(m *Options) *Message {
	return &Message{
		Options: m,
	}
}

func (m *Message) Encode() []byte {
	ba, _ := proto.Marshal(m)
	return ba
}

func ReadAndUnmarshal(conn net.Conn, call Call) {
	reader := bufio.NewReader(conn)
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		m := NewMessage(&Options{})
		err = proto.Unmarshal(buf[:n], m)
		if err != nil {
			call(nil, err)
			return
		}
		if call(m, err) {
			return
		}
	}
}

func UnmarshalUDP(data []byte) (*Message, error) {
	m := NewMessage(&Options{})
	return m, proto.Unmarshal(data, m)
}

var PreMessage1024 = make([]byte, 1024)

func init() {
	for i := range PreMessage1024 {
		PreMessage1024[i] = 'x'
	}
	PreMessage1024[len(PreMessage1024)-1] = '\n'
}
