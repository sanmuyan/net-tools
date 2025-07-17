package benchtestc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ClientConn interface {
	run(wg *sync.WaitGroup)
}

type Client struct {
	Server      string
	Protocol    string
	Timeout     time.Duration
	Interval    time.Duration
	MaxThread   int
	MaxMessages int
	ctx         context.Context
}

func NewClient(ctx context.Context, server string, protocol string, timeout int, interval int, maxThread int, maxMessages int) *Client {
	return &Client{
		Server:      server,
		Protocol:    protocol,
		Timeout:     time.Millisecond * time.Duration(timeout),
		Interval:    time.Millisecond * time.Duration(interval),
		MaxThread:   maxThread,
		MaxMessages: maxMessages,
		ctx:         ctx,
	}
}

func (c *Client) setConnDeadline(conn net.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(c.Timeout))
}

func RunClient(client *Client) {
	var clientConn ClientConn
	switch client.Protocol {
	case "tcp":
		clientConn = NewTCPClient(client)
	case "udp":
		clientConn = NewUDPClient(client)
	case "ws":
		clientConn = NewWSClient(client)
	case "http", "https":
		clientConn = NewHTTPClient(client)
	default:
		logrus.Fatalf("unknown protocol: %s", client.Protocol)
		return
	}
	wg := new(sync.WaitGroup)
	for i := 0; i < client.MaxThread; i++ {
		wg.Add(1)
		go clientConn.run(wg)
	}
	wg.Wait()
}

func Run(ctx context.Context) {
	protocol := viper.GetString("protocol")
	server := viper.GetString("server-addr")
	timeout := viper.GetInt("timeout")
	interval := viper.GetInt("interval")
	maxThread := viper.GetInt("max-thread")
	maxMessages := viper.GetInt("max-messages")
	RunClient(NewClient(ctx, server, protocol, timeout, interval, maxThread, maxMessages))
}
