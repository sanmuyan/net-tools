package speedtestc

import (
	"bufio"
	"context"
	"errors"
	"github.com/sanmuyan/xpkg/xconstant"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/speedtest"
	"sync/atomic"
	"time"
)

type TCPClient struct {
	*Client
}

func NewTCPClient(client *Client) *TCPClient {
	return &TCPClient{
		Client: client,
	}
}

func (c *TCPClient) handleDownload(ctx context.Context, conn net.Conn) {
	reader := bufio.NewReaderSize(conn, speedtest.ReadBufferSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.setConnDeadline(conn)
			msg, err := speedtest.ReadTCP(reader)
			if err != nil {
				logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
				if errors.Is(err, xconstant.BufferedTooSmallError) {
					reader.Reset(bufio.NewReaderSize(conn, speedtest.ReadBufferSize))
					continue
				}
				c.getErrorCh() <- err
				return
			}
			switch msg.GetCtl() {
			case speedtest.NewData:
				atomic.AddInt64(&c.TotalSize, int64(speedtest.TCPDataSize))
			default:
				logrus.Debugf("unknown ctl: %d %s", msg.GetCtl(), conn.RemoteAddr())
				continue
			}
		}
	}
}

func (c *TCPClient) handleUpload(ctx context.Context, conn net.Conn) {
	for {
		select {
		case <-ctx.Done():
			_ = conn.Close()
			return
		default:
			c.setConnDeadline(conn)
			err := speedtest.WriteTCP(&speedtest.Message{
				Ctl:  speedtest.NewData,
				Data: speedtest.PreMessageTCP,
			}, conn)
			if err != nil {
				logrus.Debugf("failed to write: %v %s", err, conn.RemoteAddr())
				c.getErrorCh() <- err
				return
			}
			atomic.AddInt64(&c.TotalSize, int64(speedtest.TCPDataSize))
		}
	}
}

func (c *TCPClient) run() {
	conn := c.createConn()
	defer func() {
		_ = conn.Close()
	}()
	logrus.Infof("tcp %s testing to %s", c.Mode, c.Server)
	err := speedtest.WriteTCP(&speedtest.Message{
		Ctl:      speedtest.NewTest,
		TestTime: int32(c.TestTime),
		TestMode: c.Mode,
	}, conn)
	if err != nil {
		logrus.Fatalf("failed to write to tcp server: %v", err)
	}
	ctx, cancel := context.WithTimeout(c.ctx, time.Duration(c.TestTime)*time.Second)
	defer cancel()
	switch c.Mode {
	case "download":
		c.handleDownload(ctx, conn)
	case "upload":
		c.handleUpload(ctx, conn)
	}
}

func (c *TCPClient) createConn() net.Conn {
	conn, err := net.DialTimeout("tcp", c.Server, time.Second*3)
	if err != nil {
		logrus.Fatalf("failed to dial server: %v", err)
	}
	c.setConnDeadline(conn)
	return conn
}
