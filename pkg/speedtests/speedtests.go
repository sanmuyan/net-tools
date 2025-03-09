package speedtests

import (
	"context"
	"github.com/spf13/viper"
	"net"
	"time"
)

type ServerConn interface {
	run()
}

type Server struct {
	ServerBind string
	ctx        context.Context
}

func NewServer(ctx context.Context, serverBind string) *Server {
	return &Server{
		ServerBind: serverBind,
		ctx:        ctx,
	}
}

func (s *Server) setConnDeadline(conn net.Conn, testTime int) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(testTime)))
	_ = conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(testTime)))
}

func RunServer(server *Server) {
	NewTCPServer(server).run()
}

func Run(ctx context.Context) {
	serverBind := viper.GetString("server-bind")
	RunServer(NewServer(ctx, serverBind))
}
