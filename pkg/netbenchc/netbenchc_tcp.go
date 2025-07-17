package netbenchc

import (
	"bufio"
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/netbench"
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

func (c *TCPClient) sendHandler(ctx context.Context, conn net.Conn) {
	reader := bufio.NewReaderSize(conn, netbench.ReadBufferSize)
	for {
		select {
		case <-ctx.Done():
			_ = conn.Close()
			return
		default:
			startTime := time.Now().UnixMilli()
			sendMsg := netbench.GenerateMessage(netbench.GenerateRequestID())
			logrus.Debugf("%s message: %s to %s", c.Protocol, sendMsg.GetRequestID(), conn.RemoteAddr())
			err := netbench.WriteTCP(sendMsg, conn)
			if err != nil {
				logrus.Warnf("failed to write: %v %s", err, conn.RemoteAddr())
				return
			}
			c.setConnDeadline(conn)
			receiveMsg, err := netbench.ReadTCP(reader)
			if err != nil {
				logrus.Warnf("failed to read: %v %s", err, conn.RemoteAddr())
				return
			}
			endTime := time.Now().UnixMilli()
			logrus.Infof("%s message: %sfrom %s %dms", c.Protocol, receiveMsg.GetRequestID(), conn.RemoteAddr(), endTime-startTime)
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
