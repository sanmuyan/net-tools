package speedtestc

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"time"
)

type UDPClient struct {
	*Client
}

func NewUDPClient(client *Client) *UDPClient {
	return &UDPClient{
		Client: client,
	}
}

func (c *UDPClient) run() {
	conn := c.createConn()
	c.setConnDeadline(conn)
	defer func() {
		_ = conn.Close()
	}()
	log.Printf("udp %s testing to %s", c.Mode, c.Server)
	_, err := conn.Write(c.createCtlMsg())
	if err != nil {
		logrus.Fatalf("failed to write to tcp server: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.TestTime)*time.Second)
	defer cancel()
	switch c.Mode {
	case "download":
		c.handleDownload(ctx, conn)
	case "upload":
		c.handleUpload(ctx, conn)
	}
}

func (c *UDPClient) createConn() net.Conn {
	udpAddr, err := net.ResolveUDPAddr("udp", c.Server)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		logrus.Fatalf("failed to dial server: %v", err)
	}
	return conn
}
