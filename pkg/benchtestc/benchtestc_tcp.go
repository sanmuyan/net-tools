package benchtestc

import (
	"bufio"
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/benchtest"
	"sync"
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

func (c *TCPClient) sendMessage(reader *bufio.Reader, conn net.Conn) (*int64, error) {
	startTime := time.Now().UnixMilli()
	sendMsg := benchtest.GenerateMessage(benchtest.GenerateRequestID())
	logrus.Debugf("%s message: %s to %s", c.Protocol, sendMsg.GetRequestID(), conn.RemoteAddr())
	err := benchtest.WriteTCP(sendMsg, conn)
	if err != nil {
		logrus.Warnf("failed to write: %v %s", err, conn.RemoteAddr())
		return nil, err
	}
	c.setConnDeadline(conn)
	receiveMsg, err := benchtest.ReadTCP(reader)
	if err != nil {
		logrus.Warnf("failed to read: %v %s", err, conn.RemoteAddr())
		return nil, err
	}
	timing := time.Now().UnixMilli() - startTime
	logrus.Infof("%s message: %s from %s %s", c.Protocol, receiveMsg.GetRequestID(), conn.RemoteAddr(), timeToStrUnit(timing))
	return &timing, nil
}

func (c *TCPClient) sendHandler(ctx context.Context, conn net.Conn) {
	reader := bufio.NewReaderSize(conn, benchtest.ReadBufferSize)
	for i := 0; i < c.MaxMessages || c.MaxMessages <= 0; i++ {
		select {
		case <-ctx.Done():
			_ = conn.Close()
			return
		default:
			timing, err := c.sendMessage(reader, conn)
			if err != nil {
				c.addErrorCount()
				return
			}
			c.addSuccessCount(*timing)
		}
		time.Sleep(c.Interval)
	}
}

func (c *TCPClient) run(wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := c.createConn()
	if err != nil {
		logrus.Errorf("failed to dial server: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
		logrus.Debugf("%s test finished in %s", c.Protocol, conn.RemoteAddr())
	}()
	logrus.Debugf("%s testing to %s", c.Protocol, c.Server)
	c.sendHandler(c.ctx, conn)
}

func (c *TCPClient) createConn() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", c.Server, c.Timeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
