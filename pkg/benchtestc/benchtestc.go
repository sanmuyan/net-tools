package benchtestc

import (
	"context"
	"github.com/sanmuyan/xpkg/xtime"
	"net"
	"net-tools/pkg/loger"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ClientConn interface {
	run(wg *sync.WaitGroup)
}

type Client struct {
	Server           string
	Protocol         string
	Timeout          time.Duration
	Interval         time.Duration
	MaxThread        int
	MaxMessages      int
	ctx              context.Context
	errCount         int64
	successCount     int64
	successTimeCount int64
	successMinTime   int64
	successMaxTime   int64
	mx               sync.Mutex
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
func (c *Client) addSuccessCount(t int64) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.successCount++
	c.successTimeCount += t
	if t > c.successMaxTime || c.successMaxTime == 0 {
		c.successMaxTime = t
		return
	}
	if t < c.successMinTime || c.successMinTime == 0 {
		c.successMinTime = t
		return
	}
}

func (c *Client) addErrorCount() {
	atomic.AddInt64(&c.errCount, 1)
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
	startTime := time.Now().Unix()
	wg := new(sync.WaitGroup)
	for i := 0; i < client.MaxThread; i++ {
		wg.Add(1)
		go clientConn.run(wg)
	}
	wg.Wait()
	timing := time.Now().Unix() - startTime
	var avg int64
	if client.successCount > 0 {
		avg = client.successTimeCount / client.successCount
	}
	loger.S.Infof("Success=%d, Error=%d, Timing=%ds Max=%s, Min=%s, Avg=%s",
		client.successCount, client.errCount, timing, timeToStrUnit(client.successMaxTime), timeToStrUnit(client.successMinTime), timeToStrUnit(avg))
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

func timeToStrUnit(tm int64) string {
	return xtime.TimeToStrUnitTrim(time.Duration(tm)*time.Millisecond, 3)
}
