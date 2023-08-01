package tcpping

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type TCPPing struct {
	Host     string
	Port     int
	Timeout  int
	Protocol string
}

func NewTCPPing(host string, port int, timeout int, protocol string) *TCPPing {
	if port < 0 || port > 65535 {
		port = 22
	}
	if timeout < 5 {
		timeout = 5
	}
	if timeout > 100000 {
		timeout = 100000
	}
	return &TCPPing{Host: host, Port: port, Timeout: timeout, Protocol: protocol}
}

func (t *TCPPing) httpHello(conn net.Conn) error {
	_, err := conn.Write([]byte("HEAD / HTTP/1.0\n\n"))
	if err != nil {
		return errors.New(fmt.Sprintf("http write timeout"))
	}
	_, err = conn.Read([]byte(""))
	if err != nil {
		return errors.New(fmt.Sprintf("http read timeout"))
	}
	return nil
}

func (t *TCPPing) readHello(conn net.Conn) error {
	_, err := conn.Read([]byte(""))
	if err != nil {
		return errors.New(fmt.Sprintf("ssh read timeout"))
	}
	return nil
}

func (t *TCPPing) PING(errorMessage chan string, pingTime chan int) {
	starTime := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port), time.Duration(t.Timeout)*time.Millisecond)
	if err != nil {
		errorMessage <- fmt.Sprintf("timeout")
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	pt := int(time.Now().Sub(starTime).Milliseconds())
	if pt < t.Timeout {
		// 剩余超时时间应该减去建立连接的耗时
		_t := t.Timeout - pt
		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(_t) * time.Millisecond))
		_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(_t) * time.Millisecond))
	} else {
		errorMessage <- fmt.Sprintf("timeout")
		return
	}
	switch t.Protocol {
	case "http":
		err = t.httpHello(conn)
	case "read":
		err = t.readHello(conn)
	}
	if err != nil {
		errorMessage <- fmt.Sprintf("%s", err)
		return
	}
	pt = int(time.Now().Sub(starTime).Milliseconds())
	if pt > t.Timeout {
		errorMessage <- fmt.Sprintf("timeout")
		return
	}
	pingTime <- pt
}
