package nettestc

import (
	"context"
	"github.com/sanmuyan/xpkg/xutil"
	"github.com/sirupsen/logrus"
	"net"
	"net-tools/pkg/nettest"
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

func (c *UDPClient) sendHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			startTime := time.Now().UnixMilli()
			sendMsg := nettest.GenerateMessage(nettest.GenerateRequestID())
			logrus.Debugf("udp message: %s to %s", sendMsg.GetRequestID(), c.conn.RemoteAddr())
			_, err := c.conn.Write(xutil.RemoveError(nettest.Marshal(sendMsg)))
			if err != nil {
				logrus.Warnf("failed to write: %s %s", err, c.conn.RemoteAddr())
				return
			}
			c.setConnDeadline(c.conn)
			data := make([]byte, nettest.ReadBufferSize)
			n, err := c.conn.Read(data)
			if err != nil {
				logrus.Warnf("failed to read: %s %s", err, c.conn.RemoteAddr())
				return
			}
			receiveMsg, err := nettest.Unmarshal(data[:n])
			if err != nil {
				logrus.Warnf("failed to unmarshal: %s %s", err, c.conn.RemoteAddr())
				return
			}
			endTime := time.Now().UnixMilli()
			logrus.Infof("udp message: %s from %s %dms", receiveMsg.GetRequestID(), c.conn.RemoteAddr(), endTime-startTime)
		}
		time.Sleep(c.Interval)
	}
}

func (c *UDPClient) run(wg *sync.WaitGroup) {
	defer wg.Done()
	c.conn = c.createConn()
	defer func() {
		_ = c.conn.Close()
		logrus.Debugf("udp test finished in %s", c.conn.RemoteAddr())
	}()
	c.sendHandler(c.ctx)
}

func (c *UDPClient) createConn() net.Conn {
	conn, err := net.DialTimeout("udp", c.Server, c.Timeout)
	if err != nil {
		logrus.Fatalf("failed to dial server: %v", err)
	}
	c.setConnDeadline(conn)
	return conn
}
