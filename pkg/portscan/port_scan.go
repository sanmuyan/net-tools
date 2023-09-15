package portscan

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type PortScan struct {
	ports     []int
	IP        string
	MaxThread int
	Timeout   int
}

func NewPortScan(ports []int, ip string, maxThread int, timeout int) *PortScan {
	if maxThread < 1 {
		maxThread = 1
	}
	// 最大并发不宜太多，过多并发很容易被防火墙拦截
	if maxThread > 100 {
		maxThread = 100
	}
	if timeout < 10 {
		timeout = 10
	}
	if timeout > 10000 {
		timeout = 10000
	}
	return &PortScan{
		ports:     ports,
		IP:        ip,
		MaxThread: maxThread,
		Timeout:   timeout,
	}
}

func (p *PortScan) Scan(done chan bool, openPorts chan int) {
	defer func() {
		done <- true
	}()
	maxThread := make(chan struct{}, p.MaxThread)
	wg := sync.WaitGroup{}
	for _, port := range p.ports {
		wg.Add(1)
		maxThread <- struct{}{}
		go func(port int) {
			defer func() {
				wg.Done()
				<-maxThread
			}()
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.IP, port), time.Duration(p.Timeout)*time.Millisecond)
			if err != nil {
				return
			}
			_ = conn.Close()
			openPorts <- port
		}(port)
	}
	wg.Wait()
}
