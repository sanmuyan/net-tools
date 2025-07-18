package benchtestc

import (
	"context"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/benchtest"
	"sync"
	"time"
)

type UDPClient struct {
	*Client
	conn net.Conn
}

func NewUDPClient(client *Client) *UDPClient {
	return &UDPClient{
		Client: client,
	}
}

func (c *UDPClient) sendMessage() (*int64, error) {
	startTime := time.Now().UnixMilli()
	sendMsg := benchtest.GenerateMessage(benchtest.GenerateRequestID())
	logrus.Debugf("%s message: %s to %s", c.Protocol, sendMsg.GetRequestID(), c.conn.RemoteAddr())
	_, err := c.conn.Write(xutil.RemoveError(benchtest.Marshal(sendMsg)))
	if err != nil {
		logrus.Warnf("failed to write: %s %s", err, c.conn.RemoteAddr())
		return nil, err
	}
	c.setConnDeadline(c.conn)
	data := make([]byte, benchtest.ReadBufferSize)
	n, err := c.conn.Read(data)
	if err != nil {
		logrus.Warnf("failed to read: %s %s", err, c.conn.RemoteAddr())
		return nil, err
	}
	receiveMsg, err := benchtest.Unmarshal(data[:n])
	if err != nil {
		logrus.Warnf("failed to unmarshal: %s %s", err, c.conn.RemoteAddr())
		return nil, err
	}
	timing := time.Now().UnixMilli() - startTime
	logrus.Infof("%s message: %s from %s %s", c.Protocol, receiveMsg.GetRequestID(), c.conn.RemoteAddr(), timeToStrUnit(timing))
	return &timing, err
}

func (c *UDPClient) sendHandler(ctx context.Context) {
	for i := 0; i < c.MaxMessages || c.MaxMessages <= 0; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			timing, err := c.sendMessage()
			if err != nil {
				c.addErrorCount()
				return
			}
			c.addSuccessCount(*timing)
		}
		time.Sleep(c.Interval)
	}
}

func (c *UDPClient) run(wg *sync.WaitGroup) {
	defer wg.Done()
	c.conn = c.createConn()
	defer func() {
		_ = c.conn.Close()
	}()
	c.sendHandler(c.ctx)
}

func (c *UDPClient) createConn() net.Conn {
	conn, err := net.DialTimeout("udp", c.Server, c.Timeout)
	if err != nil {
		logrus.Fatalf("failed to dial server: %v", err)
	}
	return conn
}
