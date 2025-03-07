package speedtestc

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
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

func (c *TCPClient) run() {
	conn := c.createConn()
	c.setConnDeadline(conn)
	defer func() {
		_ = conn.Close()
	}()
	logrus.Infof("tcp %s testing to %s", c.Mode, c.Server)

	_, err := conn.Write(c.createCtlMsg())
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
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Server)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logrus.Fatalf("failed to dial server: %v", err)
	}
	return conn
}
