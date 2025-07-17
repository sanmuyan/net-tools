package nettests

import (
	"bufio"
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/nettest"
)

type TCPServer struct {
	*Server
}

func NewTCPServer(server *Server) *TCPServer {
	return &TCPServer{Server: server}
}

func (s *TCPServer) replyHandler(ctx context.Context, conn net.Conn) {
	defer func() {
		_ = conn.Close()
		logrus.Debugf("tcp test finished in %s", conn.RemoteAddr())
	}()
	reader := bufio.NewReaderSize(conn, nettest.ReadBufferSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.setConnDeadline(conn)
			receiveMsg, err := nettest.ReadTCP(reader)
			if err != nil {
				logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
				return
			}
			logrus.Infof("tcp message: %s from %s", receiveMsg.GetRequestID(), conn.RemoteAddr())
			sendMsg := nettest.GenerateMessage(receiveMsg.GetRequestID())
			logrus.Debugf("tcp message: %s to %s", sendMsg.GetRequestID(), conn.RemoteAddr())
			err = nettest.WriteTCP(sendMsg, conn)
			if err != nil {
				logrus.Warnf("failed to write: %v %s", err, conn.RemoteAddr())
				return
			}
		}
	}
}

func (s *TCPServer) run() {
	listener, err := net.Listen("tcp", s.ServerBind)
	if err != nil {
		logrus.Fatalf("listen error: %v", err)
	}
	defer func() {
		_ = listener.Close()
	}()
	logrus.Infof("tcp server listening on %s", s.ServerBind)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go s.replyHandler(s.ctx, conn)
		}
	}()
	<-s.ctx.Done()
}
