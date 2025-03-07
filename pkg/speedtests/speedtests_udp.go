package speedtests

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"net-tools/pkg/speedtest"
	"sync"
	"time"
)

type updPool struct {
	totalSize int
	mu        sync.Mutex
}

type UDPServer struct {
	*Server
	conn       *net.UDPConn
	uploadPool sync.Map
}

func NewUDPServer(server *Server) *UDPServer {
	return &UDPServer{
		Server: server,
	}
}

func (s *UDPServer) handleDownload(ctx context.Context, addr *net.UDPAddr) {
	for {
		select {
		case <-ctx.Done():
			logrus.Infof("udp download finished in %s", addr.String())
			return
		default:
			_, err := s.conn.WriteTo(speedtest.PreMessage1024, addr)
			if err != nil {
				return
			}
		}
	}
}

func (s *UDPServer) handleUpload(addr *net.UDPAddr, n int) {
	if mp, ok := s.uploadPool.Load(addr.String()); ok {
		_mp := mp.(*updPool)
		_mp.mu.Lock()
		_mp.totalSize += n
		s.uploadPool.Store(addr.String(), _mp)
		_mp.mu.Unlock()
	}
}

func (s *UDPServer) controller(addr *net.UDPAddr, msg *speedtest.Message) {
	log.Printf("udp %s from %s", msg.GetCtl(), addr.String())
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*time.Duration(msg.TestTime))
	defer cancel()
	switch msg.Ctl {
	case "download":
		s.handleDownload(ctx, addr)
	case "upload":
		s.uploadPool.Store(addr.String(), &updPool{
			totalSize: 0,
		})
		time.Sleep(time.Second * time.Duration(msg.TestTime))
		logrus.Infof("upload finished in %s", addr.String())
		if up, ok := s.uploadPool.Load(addr.String()); ok {
			// 执行结束后，把客户端上传的数据总和统计返回给客户端
			_, _ = s.conn.WriteTo(speedtest.NewMessage(&speedtest.Options{
				TotalSize: int64(up.(*updPool).totalSize),
			}).Encode(), addr)
			s.uploadPool.Delete(addr.String())
		}
	}
}

func (s *UDPServer) run() {
	updAddr, err := net.ResolveUDPAddr("udp", s.ServerBind)
	s.conn, err = net.ListenUDP("udp", updAddr)
	if err != nil {
		logrus.Fatalf("listen error: %v", err)
	}
	defer func() {
		_ = s.conn.Close()
	}()
	logrus.Infof("udp server listening on %s", s.ServerBind)
	go func() {
		for {
			data := make([]byte, 1024)
			n, addr, err := s.conn.ReadFromUDP(data)
			if err != nil {
				return
			}
			// 判断客户端是否已经在连接池中，如果是统计客户端上传的数据总和
			if _, ok := s.uploadPool.Load(addr.String()); ok {
				go s.handleUpload(addr, n)
				continue
			}
			if n == 1024 {
				continue
			}
			msg, err := speedtest.UnmarshalUDP(data[:n])
			if err != nil {
				logrus.Errorf("decode udp error: %v", err)
				continue
			}
			go s.controller(addr, msg)
		}
	}()
	<-s.ctx.Done()
}
