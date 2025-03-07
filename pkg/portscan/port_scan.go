package portscan

import (
	"context"
	"fmt"
	"github.com/sanmuyan/xpkg/xnet"
	"github.com/spf13/viper"
	"net"
	"net-tools/pkg/loger"
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

func Run(ctx context.Context, args []string) {
	maxThread := viper.GetInt("max-thread")
	timeout := viper.GetInt("timeout")
	var ip string
	var port string

	if len(args) == 2 {
		ip = args[0]
		port = args[1]
	}
	if !xnet.IsIP(ip) && !xnet.IsCIDR(ip) && !xnet.IsIPRange(ip) {
		loger.S.Fatalf("Invalid ip: %v", ip)
	}
	if !xnet.IsPort(port) {
		loger.S.Fatalf("Invalid port: %v", port)
	}

	loger.S.Debugf("PortScan: %s %s", ip, port)

	ips := xnet.ParseIPList(ip)
	ports := xnet.GeneratePorts(port)
	for _, ip := range ips {
		done := make(chan bool)
		openPorts := make(chan int)
		p := NewPortScan(ports, ip, maxThread, timeout)
		go p.Scan(done, openPorts)
		func() {
			for {
				select {
				case openPort := <-openPorts:
					loger.S.Infof("Opened: %s:%d", ip, openPort)
				case <-done:
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}
