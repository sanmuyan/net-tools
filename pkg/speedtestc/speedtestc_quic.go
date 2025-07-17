package speedtestc

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"github.com/sirupsen/logrus"
	"net-tools/pkg/speedtest"
	"sync/atomic"
	"time"
)

type QUICClient struct {
	*Client
}

func NewQUICClient(client *Client) *QUICClient {
	return &QUICClient{
		Client: client,
	}
}

func (c *QUICClient) handleDownload(ctx context.Context, conn *quic.Conn, stream *quic.Stream) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.setConnDeadline(stream)
			msg, err := speedtest.ReadQUIC(stream)
			if err != nil {
				logrus.Debugf("failed to read: %v %s", err, conn.RemoteAddr())
				c.getErrorCh() <- err
				return
			}
			switch msg.GetCtl() {
			case speedtest.NewData:
				atomic.AddInt64(&c.TotalSize, int64(speedtest.QUICDataSize))
			default:
				logrus.Debugf("unknown ctl: %d %s", msg.GetCtl(), conn.RemoteAddr())
				continue
			}
		}
	}
}

func (c *QUICClient) handleUpload(ctx context.Context, conn *quic.Conn, stream *quic.Stream) {
	for {
		select {
		case <-ctx.Done():
			_ = stream.Close()
			return
		default:
			c.setConnDeadline(stream)
			err := speedtest.WriteQUIC(&speedtest.StMessage{
				Ctl:  speedtest.NewData,
				Data: speedtest.PreMessageQUIC,
			}, stream)
			if err != nil {
				logrus.Debugf("failed to write: %v %s", err, conn.RemoteAddr())
				c.getErrorCh() <- err
				return
			}
			atomic.AddInt64(&c.TotalSize, int64(speedtest.QUICDataSize))
		}
	}
}

func (c *QUICClient) setConnDeadline(stream *quic.Stream) {
	_ = stream.SetReadDeadline(time.Now().Add(time.Second * 3))
	_ = stream.SetWriteDeadline(time.Now().Add(time.Second * 3))
}

func (c *QUICClient) run() {
	conn := c.createConn()
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		logrus.Fatalf("failed to open stream: %v", err)
	}
	defer func() {
		_ = stream.Close()
	}()
	logrus.Infof("quic %s testing to %s", c.Mode, c.Server)
	err = speedtest.WriteQUIC(&speedtest.StMessage{
		Ctl:      speedtest.NewTest,
		TestTime: int32(c.TestTime),
		TestMode: c.Mode,
	}, stream)
	if err != nil {
		logrus.Fatalf("failed to write to quic server: %v", err)
	}
	ctx, cancel := context.WithTimeout(c.ctx, time.Duration(c.TestTime)*time.Second)
	defer cancel()
	switch c.Mode {
	case "download":
		c.handleDownload(ctx, conn, stream)
	case "upload":
		c.handleUpload(ctx, conn, stream)
	}
}

func (c *QUICClient) createConn() *quic.Conn {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	quicConfig := &quic.Config{}
	conn, err := quic.DialAddr(context.Background(), c.Server, tlsConfig, quicConfig)
	if err != nil {
		logrus.Fatalf("failed to dial server: %v", err)
	}
	return conn
}
