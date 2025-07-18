package benchtests

import (
	"context"
	"github.com/quic-go/quic-go"
	"github.com/sanmuyan/xpkg/xcrypto"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"net-tools/pkg/benchtest"
	"time"
)

type QUICServer struct {
	*Server
}

func NewQUICServer(server *Server) *QUICServer {
	return &QUICServer{Server: server}
}

func (s *QUICServer) replyHandler(ctx context.Context, conn *quic.Conn) {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		logrus.Debugf("failed to accept stream: %v", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.setConnDeadline(stream)
			receiveMsg, err := benchtest.ReadQUIC(stream)
			if err != nil {
				logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
				return
			}
			logrus.Infof("%s message: %s from %s", s.Protocol, receiveMsg.GetRequestID(), conn.RemoteAddr())
			sendMsg := benchtest.GenerateMessage(receiveMsg.GetRequestID())
			logrus.Debugf("%s message: %s to %s", s.Protocol, sendMsg.GetRequestID(), conn.RemoteAddr())
			err = benchtest.WriteQUIC(sendMsg, stream)
			if err != nil {
				logrus.Warnf("failed to write: %v %s", err, conn.RemoteAddr())
				return
			}
		}
	}
}

func (s *QUICServer) run() {
	tlsConfig := xutil.RemoveError(xcrypto.CreateCertToTLS(nil))
	quicConfig := &quic.Config{}
	listener, err := quic.ListenAddr(s.ServerBind, tlsConfig, quicConfig)
	if err != nil {
		logrus.Fatalf("listen error: %v", err)
	}
	defer func() {
		_ = listener.Close()
	}()
	logrus.Infof("%s server listening on %s", s.Protocol, s.ServerBind)
	go func() {
		for {
			conn, err := listener.Accept(s.ctx)
			if err != nil {
				continue
			}
			go s.replyHandler(s.ctx, conn)
		}
	}()
	<-s.ctx.Done()
}

func (s *QUICServer) setConnDeadline(stream *quic.Stream) {
	_ = stream.SetReadDeadline(time.Now().Add(s.Timeout))
}
