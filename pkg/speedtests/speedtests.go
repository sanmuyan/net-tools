package speedtests

import (
	"context"
	"github.com/spf13/viper"
)

type ServerConn interface {
	run()
}

type Server struct {
	ServerBind string
	Protocol   string
	ctx        context.Context
}

func NewServer(ctx context.Context, serverBind string, protocol string) *Server {
	return &Server{
		ServerBind: serverBind,
		Protocol:   protocol,
		ctx:        ctx,
	}
}

func RunServer(server *Server) {
	//switch server.Protocol {
	//case "tcp":
	//	NewTCPServer(server).run()
	//case "udp":
	//	NewUDPServer(server).run()
	//default:
	//	go NewUDPServer(server).run()
	//	NewTCPServer(server).run()
	//}
	NewTCPServer(server).run()
}

func Run(ctx context.Context) {
	serverBind := viper.GetString("server-bind")
	protocol := viper.GetString("protocol")
	RunServer(NewServer(ctx, serverBind, protocol))
}
