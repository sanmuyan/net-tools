package speedtestc

import (
	"bufio"
	"context"
	"github.com/sanmuyan/xpkg/xnet"
	"log"
	"net"
	"net-tools/pkg/speedtest"
	"sync"
	"time"
)

type ClientConn interface {
	run()
	getTotalSize() int
}

type Client struct {
	ServerHost string
	ServerPort int
	Mode       string
	TestTime   int
	Protocol   string
	MaxThread  int
	TotalSize  int
	mu         sync.Mutex
}

func NewClient(serverHost string, serverPort int, mode string, testTime int, protocol string, maxThread int) *Client {
	return &Client{
		ServerHost: serverHost,
		ServerPort: serverPort,
		Mode:       mode,
		TestTime:   testTime,
		Protocol:   protocol,
		MaxThread:  maxThread,
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

func Start(client *Client) {
	var c ClientConn
	if client.Protocol == "tcp" {
		c = NewTCPClient(client)
	} else {
		c = NewUDPClient(client)
	}
	t := time.NewTicker(time.Second)
	defer t.Stop()
	go func() {
		var latestSize int
		for range t.C {
			log.Printf("real-time speed: %s", xnet.GetDataSpeed(c.getTotalSize()-latestSize, 1))
			latestSize = c.getTotalSize()
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
	log.Printf("finished avg speed: %s", xnet.GetDataSpeed(c.getTotalSize(), client.TestTime))
}
