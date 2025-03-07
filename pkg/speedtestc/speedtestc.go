package speedtestc

import (
	"bufio"
	"context"
	"github.com/sanmuyan/xpkg/xnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"net-tools/pkg/speedtest"
	"sync"
	"time"
)

type ClientConn interface {
	run()
	getTotalSize() int
	getErrorCh() chan error
}

type Client struct {
	Server    string
	Mode      string
	TestTime  int
	Protocol  string
	MaxThread int
	TotalSize int
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

func (c *Client) handleDownload(ctx context.Context, conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := reader.ReadBytes('\n')
			if err != nil {
				c.getErrorCh() <- err
				return
			}
			go func() {
				c.mu.Lock()
				c.TotalSize += len(data)
				c.mu.Unlock()
			}()
		}
	}
}

func (c *Client) handleUpload(ctx context.Context, conn net.Conn) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, err := conn.Write(speedtest.PreMessage1024)
				if err != nil {
					c.getErrorCh() <- err
					return
				}
			}
		}
	}()
	// 等待服务端发送上传数据总和
	speedtest.ReadAndUnmarshal(conn, func(msg *speedtest.Message, err error) (exit bool) {
		if err != nil {
			return true
		}
		c.mu.Lock()
		c.TotalSize += int(msg.GetTotalSize())
		c.mu.Unlock()
		return true
	})
}

func (c *Client) setConnDeadline(conn net.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.TestTime+5)))
	_ = conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(c.TestTime+5)))
}

func (c *Client) getTotalSize() int {
	return c.TotalSize
}

func (c *Client) createCtlMsg() []byte {
	return speedtest.NewMessage(&speedtest.Options{
		Ctl:      c.Mode,
		TestTime: int32(c.TestTime),
	}).Encode()

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
	if client.Protocol == "tcp" {
		c = NewTCPClient(client)
	} else {
		c = NewUDPClient(client)
	}
	var runTime int
	t := time.NewTicker(time.Second)
	defer t.Stop()
	go func() {
		var latestSize int
		for range t.C {
			logrus.Infof("real-time speed: %s", xnet.GetDataSpeed(c.getTotalSize()-latestSize, 1))
			latestSize = c.getTotalSize()
			runTime += 1
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
		logrus.Infof("finished avg speed: %s", xnet.GetDataSpeed(c.getTotalSize(), testTime))
	}
}
