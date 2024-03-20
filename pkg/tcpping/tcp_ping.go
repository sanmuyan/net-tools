package tcpping

import (
	"crypto/tls"
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

var (
	writeError = errors.New(fmt.Sprintf("write error"))
	readError  = errors.New(fmt.Sprintf("read timeout"))
	timeoutMsg = fmt.Sprintf("timeout")
)

func (t *TCPPing) httpHelloMessage() []byte {
	return []byte("HEAD / HTTP/1.1\nHost: " + t.Host + "\n\n")
}

func (t *TCPPing) httpHello(conn net.Conn) error {
	_, err := conn.Write(t.httpHelloMessage())
	if err != nil {
		return writeError
	}
	_, err = conn.Read([]byte(""))
	if err != nil {
		return readError
	}
	return nil
}

func (t *TCPPing) readHello(conn net.Conn) error {
	_, err := conn.Read([]byte(""))
	if err != nil {
		return readError
	}
	return nil
}

func (t *TCPPing) httpsHello(conn net.Conn) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	tlsConn := tls.Client(conn, tlsConfig)
	_, err := tlsConn.Write(t.httpHelloMessage())
	if err != nil {
		return writeError
	}
	_, err = conn.Read([]byte(""))
	if err != nil {
		return readError
	}
	return nil
}

func (t *TCPPing) PING(errorMessage chan string, pingTime chan int) {
	starTime := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port), time.Duration(t.Timeout)*time.Millisecond)
	if err != nil {
		errorMessage <- timeoutMsg
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
		errorMessage <- timeoutMsg
		return
	}
	switch t.Protocol {
	case "http":
		err = t.httpHello(conn)
	case "read":
		err = t.readHello(conn)
	case "https":
		err = t.httpsHello(conn)
	}
	if err != nil {
		errorMessage <- fmt.Sprintf("%s", err)
		return
	}
	pt = int(time.Now().Sub(starTime).Milliseconds())
	if pt > t.Timeout {
		errorMessage <- timeoutMsg
		return
	}
	pingTime <- pt
}
