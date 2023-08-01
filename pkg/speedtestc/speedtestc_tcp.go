package speedtestc

import (
	"context"
	"log"
	"net"
	"strconv"
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
	log.Printf("tcp %s testing to %s:%d", c.Mode, c.ServerHost, c.ServerPort)

	_, err := conn.Write(c.createCtlMsg())
	if err != nil {
		log.Fatalf("failed to write to tcp server: %v", err)
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

func (c *TCPClient) createConn() net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(c.ServerHost, strconv.Itoa(c.ServerPort)))
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	return conn
}
