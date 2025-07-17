package speedtests

import (
	"context"
	"github.com/quic-go/quic-go"
	"github.com/sanmuyan/xpkg/xcrypto"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"net-tools/pkg/speedtest"
	"time"
)

type QUICServer struct {
	*Server
}

func NewQUICServer(server *Server) *QUICServer {
	return &QUICServer{
		Server: server,
	}
}
func (s *QUICServer) handleDownload(ctx context.Context, conn *quic.Conn, stream *quic.Stream) {
	defer func() {
		logrus.Infof("quic download finished in %s", conn.RemoteAddr())
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := speedtest.WriteQUIC(&speedtest.StMessage{
				Ctl:  speedtest.NewData,
				Data: speedtest.PreMessageQUIC,
			}, stream)
			if err != nil {
				logrus.Debugf("failed to wite: %v %s", err, conn.RemoteAddr())
				return
			}
		}
	}
}

func (s *QUICServer) handleUpload(ctx context.Context, conn *quic.Conn, stream *quic.Stream) {
	defer func() {
		logrus.Infof("quic upload finished in %s", conn.RemoteAddr())
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := speedtest.ReadQUIC(stream)
			if err != nil {
				logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
				return
			}
			switch msg.GetCtl() {
			case speedtest.NewData:
				continue
			default:
				logrus.Debugf("unknown ctl: %d %s", msg.GetCtl(), conn.RemoteAddr())
				continue
			}
		}
	}
}

func (s *QUICServer) controller(conn *quic.Conn) {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		logrus.Debugf("failed to accept stream: %v", err)
		return
	}
	defer func() {
		_ = stream.Close()
	}()
	msg, err := speedtest.ReadQUIC(stream)
	if err != nil {
		logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
		return
	}
	if msg.GetCtl() != speedtest.NewTest {
		logrus.Debugf("failed ctl: %d %s", msg.GetCtl(), conn.RemoteAddr())
		return
	}
	logrus.Infof("quic %s from %s", msg.GetTestMode(), conn.RemoteAddr())
	s.setConnDeadline(stream, int(msg.GetTestTime()+1))
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*time.Duration(msg.TestTime+1))
	defer cancel()
	switch msg.GetTestMode() {
	case "download":
		s.handleDownload(ctx, conn, stream)
	case "upload":
		s.handleUpload(ctx, conn, stream)
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
	logrus.Infof("quic server listening on %s", s.ServerBind)
	go func() {
		for {
			conn, err := listener.Accept(s.ctx)
			if err != nil {
				continue
			}
			go s.controller(conn)
		}
	}()
	<-s.ctx.Done()
}

func (s *QUICServer) setConnDeadline(stream *quic.Stream, testTime int) {
	_ = stream.SetReadDeadline(time.Now().Add(time.Second * time.Duration(testTime)))
	_ = stream.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(testTime)))
}
