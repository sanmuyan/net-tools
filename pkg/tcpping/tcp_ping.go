package tcpping

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/sanmuyan/xpkg/xnet"
	"github.com/spf13/viper"
	"net"
	"net-tools/pkg/loger"
	"strconv"
	"time"
)

type TCPPing struct {
	Host     string
	Port     int
	Timeout  int64
	Protocol string
	WithTLS  bool
}

func NewTCPPing(host string, port int, timeout int64, protocol string, withTLS bool) *TCPPing {
	if port < 0 || port > 65535 {
		port = 22
	}
	if timeout < 5 {
		timeout = 5
	}
	if timeout > 100000 {
		timeout = 100000
	}
	// timeout * 1000 程序内部使用微妙
	return &TCPPing{Host: host, Port: port, Timeout: timeout * 1000, Protocol: protocol, WithTLS: withTLS}
}

var (
	timeoutMsg = errors.New(fmt.Sprintf("timeout"))
)

func (t *TCPPing) httpHelloMessage() []byte {
	return []byte("HEAD / HTTP/1.1\nHost: " + t.Host + "\n\n")
}

func (t *TCPPing) httpHello(conn net.Conn) error {
	_, err := conn.Write(t.httpHelloMessage())
	if err != nil {
		return err
	}
	_, err = conn.Read([]byte(""))
	if err != nil {
		return err
	}
	return nil
}

func (t *TCPPing) httpsHello(conn net.Conn) error {
	_, err := conn.Write(t.httpHelloMessage())
	if err != nil {
		return err
	}
	_, err = conn.Read([]byte(""))
	if err != nil {
		return err
	}
	return nil
}

func (t *TCPPing) PING(errorMessage chan error, pingTime chan float32) {
	starTime := time.Now().UnixMicro()
	var err error
	var conn net.Conn
	if t.WithTLS {
		var tlsConn *tls.Conn
		tlsConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port), &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			errorMessage <- err
			return
		}
		conn = tlsConn.NetConn()
	} else {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port), time.Duration(t.Timeout)*time.Microsecond)
		if err != nil {
			errorMessage <- err
			return
		}
	}
	defer func() {
		_ = conn.Close()
	}()
	pt := time.Now().UnixMicro() - starTime
	if pt < t.Timeout {
		// 剩余超时时间应该减去建立连接的耗时
		_t := t.Timeout - pt
		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(_t) * time.Microsecond))
		_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(_t) * time.Microsecond))
	} else {
		errorMessage <- timeoutMsg
		return
	}
	switch t.Protocol {
	case "http":
		err = t.httpHello(conn)
	default:
	}
	if err != nil {
		errorMessage <- err
		return
	}
	pt = time.Now().Add(time.Microsecond).UnixMicro() - starTime
	if pt > t.Timeout {
		errorMessage <- timeoutMsg
		return
	}
	pingTime <- float32(pt)
}

func Run(ctx context.Context, args []string) {
	protocol := viper.GetString("protocol")
	timeout := viper.GetInt("timeout")
	count := viper.GetInt("count")
	interval := viper.GetInt("interval")
	withTLS := viper.GetBool("tls")

	var host string
	var port string
	if len(args) >= 1 {
		host = args[0]
	}
	if !xnet.IsIP(host) {
		_, err := net.LookupHost(host)
		if err != nil {
			loger.S.Fatalf("ping: Name or service not known %s", host)
		}
	}
	if len(args) >= 2 {
		port = args[1]
	}
	if !xnet.IsPort(port) || port == "" {
		loger.S.Fatalf("ping: Invalid port %s", port)
	}
	portInt, _ := strconv.Atoi(port)
	p := NewTCPPing(host, portInt, int64(timeout), protocol, withTLS)
	errorMessage := make(chan error)
	pingTime := make(chan float32)
	go func() {
		for i := 0; i < count; i++ {
			p.PING(errorMessage, pingTime)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()
	var totalTime float32
	var successTotal int
	var errorTotal int
	var maxTime float32
	var minTime float32
	for i := 0; i < count; i++ {
		select {
		case m := <-errorMessage:
			loger.S.Infof("Reply from %s:%d error=%s", host, portInt, m)
			errorTotal++
		case t := <-pingTime:
			t = t / 1000
			loger.S.Infof("Reply from %s:%d time=%.3fms", host, portInt, t)
			totalTime += t
			successTotal++
			if t > maxTime {
				maxTime = t
			}
			if t < minTime || minTime == 0 {
				minTime = t
			}
		case <-ctx.Done():

		}
	}
	var avg float32
	if successTotal > 0 {
		avg = totalTime / float32(successTotal)
	}
	loger.S.Infof("Success=%d, Error=%d, Max=%.3fms, Min=%.3fms, Avg=%.3fms", successTotal, errorTotal, maxTime, minTime, avg)
}
