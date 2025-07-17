package speedtestc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/sanmuyan/xpkg/xnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ClientConn interface {
	run()
	getTotalSize() int64
	getErrorCh() chan error
}

type Client struct {
	Server    string
	Mode      string
	TestTime  int
	Protocol  string
	MaxThread int
	TotalSize int64
	mu        sync.Mutex
	ctx       context.Context
	errCh     chan error
}

func NewClient(ctx context.Context, server string, mode string, testTime int, protocol string, maxThread int) *Client {
	return &Client{
		Server:    server,
		Mode:      mode,
		TestTime:  testTime,
		Protocol:  protocol,
		MaxThread: maxThread,
		ctx:       ctx,
		errCh:     make(chan error, maxThread),
	}
}

func (c *Client) setConnDeadline(conn net.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
}

func (c *Client) getTotalSize() int64 {
	return c.TotalSize
}

func (c *Client) getErrorCh() chan error {
	return c.errCh
}

func Run(ctx context.Context) {
	testTime := viper.GetInt("test-time")
	testMode := viper.GetString("test-mode")
	protocol := viper.GetString("protocol")
	maxThread := viper.GetInt("max-thread")
	server := viper.GetString("server-addr")
	client := NewClient(ctx, server, testMode, testTime, protocol, maxThread)

	var c ClientConn
	switch protocol {
	case "tcp":
		c = NewTCPClient(client)
	case "quic":
		c = NewQUICClient(client)
	default:
		logrus.Fatalf("unknown protocol: %s", protocol)
	}
	go func() {
		var latestSize int64
		for i := 0; i < client.TestTime; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second)
				if c.getTotalSize() == 0 || c.getTotalSize() == latestSize {
					continue
				}
				logrus.Infof("real-time speed: %s", xnet.GetDataSpeed(int(c.getTotalSize()-latestSize), 1))
				latestSize = c.getTotalSize()
			}
		}
	}()
	wg := new(sync.WaitGroup)
	for i := 0; i < client.MaxThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.run()
		}()
	}
	wg.Wait()
	select {
	case <-ctx.Done():
		return
	case err := <-c.getErrorCh():
		logrus.Errorf("error: %v", err)
		return
	default:
		logrus.Infof("finished avg speed: %s", xnet.GetDataSpeed(int(c.getTotalSize()), testTime))
	}
}
