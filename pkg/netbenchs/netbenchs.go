package netbenchs

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"time"
)

type ServerConn interface {
	run()
}

type Server struct {
	ServerBind string
	Protocol   string
	Timeout    time.Duration
	ctx        context.Context
}

func NewServer(ctx context.Context, serverBind string, protocol string, timeout uint32) *Server {
	return &Server{
		ServerBind: serverBind,
		Protocol:   protocol,
		Timeout:    time.Millisecond * time.Duration(timeout),
		ctx:        ctx,
	}
}

func (s *Server) setConnDeadline(conn net.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))
}

func RunServer(server *Server) {
	var serverConn ServerConn
	switch server.Protocol {
	case "tcp":
		serverConn = NewTCPServer(server)
	case "udp":
		serverConn = NewUDPServer(server)
	case "http", "ws", "https":
		serverConn = NewHTTPServer(server)
	default:
		logrus.Fatalf("unknown protocol: %s", server.Protocol)
		return
	}
	serverConn.run()
}

func Run(ctx context.Context) {
	serverBind := viper.GetString("server-bind")
	protocol := viper.GetString("protocol")
	timeout := viper.GetUint32("timeout")
	RunServer(NewServer(ctx, serverBind, protocol, timeout))
}
