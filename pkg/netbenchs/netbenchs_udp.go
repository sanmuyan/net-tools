package netbenchs

import (
	"context"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/netbench"
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
		receiveMsg, err := netbench.Unmarshal(data)
		if err != nil {
			logrus.Warnf("failed to unmarshal message: %s", err)
			return
		}
		logrus.Infof("%s message: %s from %s", s.Protocol, receiveMsg.GetRequestID(), addr)
		sendMsg := netbench.GenerateMessage(receiveMsg.GetRequestID())
		logrus.Debugf("%s message: %s to %s", s.Protocol, sendMsg.GetRequestID(), addr)
		_, err = s.conn.WriteTo(xutil.RemoveError(netbench.Marshal(sendMsg)), addr)
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
	logrus.Infof("%s server listening on %s", s.Protocol, s.ServerBind)
	go func() {
		for {
			data := make([]byte, netbench.ReadBufferSize)
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
