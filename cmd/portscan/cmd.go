package main

import (
	"flag"
	"fmt"
	"github.com/sanmuyan/xpkg/xnet"
	"net-tools/pkg/portscan"
	"os"
)

func main() {
	maxThread := flag.Int("t", 1, "PortScan max thread")
	timeout := flag.Int("T", 200, "PortScan Timeout (ms)")
	flag.Parse()

	if len(os.Args) < 3 {
		println("Example: portscan 192.168.1.1 22")
		os.Exit(1)
	}
	if !xnet.IsIP(os.Args[1]) && !xnet.IsCIDR(os.Args[1]) && !xnet.IsIPRange(os.Args[1]) {
		println("Invalid ip: ", os.Args[1])
		os.Exit(1)
	}
	if !xnet.IsPort(os.Args[2]) {
		println("Invalid port: ", os.Args[2])
		os.Exit(1)
	}
	ips := xnet.ParseIPList(os.Args[1])
	ports := xnet.GeneratePorts(os.Args[2])
	for _, ip := range ips {
		done := make(chan bool)
		openPorts := make(chan int)
		p := portscan.NewPortScan(ports, ip, *maxThread, *timeout)
		go p.Scan(done, openPorts)
		func() {
			for {
				select {
				case openPort := <-openPorts:
					fmt.Printf("Opened: %s:%d\n", ip, openPort)
				case <-done:
					return
				}
			}
		}()
	}
}
