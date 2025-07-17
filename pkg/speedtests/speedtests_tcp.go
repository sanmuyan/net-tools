package speedtests

import (
	"bufio"
	"context"
	"errors"
	"github.com/sanmuyan/xpkg/xconstant"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/speedtest"
	"time"
)

type TCPServer struct {
	*Server
}

func NewTCPServer(server *Server) *TCPServer {
	return &TCPServer{
		Server: server,
	}
}

func (s *TCPServer) handleDownload(ctx context.Context, conn net.Conn) {
	defer func() {
		logrus.Infof("tcp download finished in %s", conn.RemoteAddr())
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := speedtest.WriteTCP(&speedtest.StMessage{
				Ctl:  speedtest.NewData,
				Data: speedtest.PreMessageTCP,
			}, conn)
			if err != nil {
				logrus.Debugf("failed to wite: %v %s", err, conn.RemoteAddr())
				return
			}
		}
	}
}

func (s *TCPServer) handleUpload(ctx context.Context, conn net.Conn) {
	defer func() {
		logrus.Infof("tcp upload finished in %s", conn.RemoteAddr())
	}()
	reader := bufio.NewReaderSize(conn, speedtest.ReadBufferSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := speedtest.ReadTCP(reader)
			if err != nil {
				logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
				if errors.Is(err, xconstant.BufferedTooSmallError) {
					reader.Reset(bufio.NewReaderSize(conn, speedtest.ReadBufferSize))
					continue
				}
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

func (s *TCPServer) controller(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	reader := bufio.NewReaderSize(conn, speedtest.ReadBufferSize)
	msg, err := speedtest.ReadTCP(reader)
	if err != nil {
		logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
		return
	}
	if msg.GetCtl() != speedtest.NewTest {
		logrus.Debugf("failed ctl: %d %s", msg.GetCtl(), conn.RemoteAddr())
		return
	}
	logrus.Infof("tcp %s from %s", msg.GetTestMode(), conn.RemoteAddr())
	s.setConnDeadline(conn, int(msg.GetTestTime()+1))
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*time.Duration(msg.TestTime+1))
	defer cancel()
	switch msg.GetTestMode() {
	case "download":
		s.handleDownload(ctx, conn)
	case "upload":
		s.handleUpload(ctx, conn)
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
			go s.controller(conn)
		}
	}()
	<-s.ctx.Done()
}
