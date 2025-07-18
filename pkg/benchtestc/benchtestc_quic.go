package benchtestc

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"github.com/sirupsen/logrus"
	"net-tools/pkg/benchtest"
	"sync"
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

func (c *QUICClient) sendMessage(conn *quic.Conn, stream *quic.Stream) (*int64, error) {
	startTime := time.Now().UnixMilli()
	sendMsg := benchtest.GenerateMessage(benchtest.GenerateRequestID())
	logrus.Debugf("%s message: %s to %s", c.Protocol, sendMsg.GetRequestID(), conn.RemoteAddr())
	err := benchtest.WriteQUIC(sendMsg, stream)
	if err != nil {
		logrus.Warnf("failed to write: %v %s", err, conn.RemoteAddr())
		return nil, err
	}
	c.setConnDeadline(stream)
	receiveMsg, err := benchtest.ReadQUIC(stream)
	if err != nil {
		logrus.Warnf("failed to read: %v %s", err, conn.RemoteAddr())
		return nil, err
	}
	timing := time.Now().UnixMilli() - startTime
	logrus.Infof("%s message: %s from %s %s", c.Protocol, receiveMsg.GetRequestID(), conn.RemoteAddr(), timeToStrUnit(timing))
	return &timing, nil
}

func (c *QUICClient) sendHandler(ctx context.Context, conn *quic.Conn, stream *quic.Stream) {
	for i := 0; i < c.MaxMessages || c.MaxMessages <= 0; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			timing, err := c.sendMessage(conn, stream)
			if err != nil {
				c.addErrorCount()
				return
			}
			c.addSuccessCount(*timing)
		}
		time.Sleep(c.Interval)
	}
}

func (c *QUICClient) run(wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := c.createConn()
	if err != nil {
		logrus.Errorf("failed to dial server: %v", err)
		return
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		logrus.Errorf("failed to open stream: %v", err)
		return
	}
	defer func() {
		_ = stream.Close()
	}()
	c.sendHandler(c.ctx, conn, stream)
}

func (c *QUICClient) createConn() (*quic.Conn, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	quicConfig := &quic.Config{}
	conn, err := quic.DialAddr(context.Background(), c.Server, tlsConfig, quicConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *QUICClient) setConnDeadline(stream *quic.Stream) {
	_ = stream.SetReadDeadline(time.Now().Add(c.Timeout))
}
