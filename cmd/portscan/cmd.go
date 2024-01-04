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
	ip := flag.String("i", "127.0.0.1", "PortScan ip")
	port := flag.String("p", "80", "PortScan port")
	flag.Parse()

	if flag.NArg() == 2 {
		*ip = flag.Arg(0)
		*port = flag.Arg(1)
	}
	if !xnet.IsIP(*ip) && !xnet.IsCIDR(*ip) && !xnet.IsIPRange(*ip) {
		println("Invalid ip: ", *ip)
		println("Example: portscan 192.168.1.1 22")
		os.Exit(1)
	}
	if !xnet.IsPort(*port) {
		println("Invalid port: ", *port)
		println("Example: portscan 192.168.1.1 22")
		os.Exit(1)
	}

	fmt.Printf("PortScan: %s %s\n", *ip, *port)

	ips := xnet.ParseIPList(*ip)
	ports := xnet.GeneratePorts(*port)
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
