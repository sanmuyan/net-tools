package speedtestc

import (
	"context"
	"log"
	"net"
	"strconv"
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
	log.Printf("udp %s testing to %s:%d", c.Mode, c.ServerHost, c.ServerPort)
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

func (c *UDPClient) createConn() net.Conn {
	udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(c.ServerHost, strconv.Itoa(c.ServerPort)))
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	return conn
}
