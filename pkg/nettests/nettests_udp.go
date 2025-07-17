package nettests

import (
	"context"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/nettest"
)

type UDPServer struct {
	*Server
	conn *net.UDPConn
}

func NewUDPServer(server *Server) *UDPServer {
	return &UDPServer{Server: server}
}

func (s *UDPServer) replyHandler(ctx context.Context, addr net.Addr, data []byte) {
	select {
	case <-ctx.Done():
		return
	default:
		receiveMsg, err := nettest.Unmarshal(data)
		if err != nil {
			logrus.Warnf("failed to unmarshal message: %s", err)
			return
		}
		logrus.Infof("udp message: %s from %s", receiveMsg.GetRequestID(), addr)
		sendMsg := nettest.GenerateMessage(receiveMsg.GetRequestID())
		logrus.Debugf("udp message: %s to %s", sendMsg.GetRequestID(), addr)
		_, err = s.conn.WriteTo(xutil.RemoveError(nettest.Marshal(sendMsg)), addr)
		if err != nil {
			logrus.Warnf("failed to write: %s %s", err, s.conn.RemoteAddr())
			return
		}
	}
}

func (s *UDPServer) run() {
	var err error
	s.conn, err = net.ListenUDP("udp", xutil.RemoveError(net.ResolveUDPAddr("udp", s.ServerBind)))
	if err != nil {
		logrus.Fatalf("listen error: %v", err)
	}
	defer func() {
		_ = s.conn.Close()
	}()
	logrus.Infof("udp server listening on %s", s.ServerBind)
	go func() {
		for {
			data := make([]byte, nettest.ReadBufferSize)
			n, addr, err := s.conn.ReadFrom(data)
			if err != nil {
				logrus.Warnf("failed to read: %v %s", err, addr)
				continue
			}
			go s.replyHandler(s.ctx, addr, data[:n])
		}
	}()
	<-s.ctx.Done()
}
