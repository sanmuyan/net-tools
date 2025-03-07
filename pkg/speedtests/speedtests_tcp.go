package speedtests

import (
	"bufio"
	"context"
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

func (s *TCPServer) handleDownload(ctx context.Context, conn *net.TCPConn) {
	defer func() {
		logrus.Infof("download finished in %s", conn.RemoteAddr())
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := conn.Write(speedtest.PreMessage1024)
			if err != nil {
				return
			}
			s.setConnDeadline(conn)
		}
	}
}

func (s *TCPServer) handleUpload(ctx context.Context, conn *net.TCPConn) {
	defer func() {
		logrus.Infof("tcp upload finished in %s", conn.RemoteAddr())
	}()
	totalSize := 0
	reader := bufio.NewReader(conn)
	defer func() {
		// 执行结束后，把客户端上传的数据总和统计返回给客户端
		_, _ = conn.Write(speedtest.NewMessage(&speedtest.Options{
			TotalSize: int64(totalSize),
		}).Encode())
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}
			totalSize += len(data)
			s.setConnDeadline(conn)
		}
	}
}

func (s *TCPServer) controller(conn *net.TCPConn) {
	defer func() {
		_ = conn.Close()
	}()
	speedtest.ReadAndUnmarshal(conn, func(msg *speedtest.Message, err error) (exit bool) {
		if err != nil {
			return true
		}
		logrus.Infof("tcp %s from %s", msg.GetCtl(), conn.RemoteAddr())
		ctx, cancel := context.WithTimeout(s.ctx, time.Second*time.Duration(msg.TestTime))
		defer cancel()
		switch msg.Ctl {
		case "download":
			s.handleDownload(ctx, conn)
		case "upload":
			s.handleUpload(ctx, conn)
		}
		return true
	})
}
func (s *TCPServer) setConnDeadline(conn *net.TCPConn) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
}

func (s *TCPServer) run() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.ServerBind)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logrus.Fatalf("listen error: %v", err)
	}
	defer func() {
		_ = listener.Close()
	}()
	logrus.Infof("tcp server listening on %s", s.ServerBind)
	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				logrus.Errorf("accept error: %v", err)
				continue
			}
			go s.controller(conn)
		}
	}()
	<-s.ctx.Done()
}
