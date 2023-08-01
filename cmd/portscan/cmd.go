package main

import (
	"flag"
	"fmt"
	"net-tools/pkg/portscan"
	"net-tools/pkg/util"
)

func main() {
	ports := flag.String("p", "0-65535", "ports to scan")
	maxThread := flag.Int("t", 1, "scan max thread")
	ipStr := flag.String("i", "127.0.0.1", "scan ip or ip net")
	timeout := flag.Int("T", 200, "timeout in ms")
	flag.Parse()

	minPort, maxPort := util.ParsePorts(*ports)
	ips := util.ParseIPList(*ipStr)

	for _, ip := range ips {
		done := make(chan bool)
		openPorts := make(chan int)
		p := portscan.NewPortScan(minPort, maxPort, ip, *maxThread, *timeout)
		go p.Scan(done, openPorts)
		func() {
			for {
				select {
				case openPort := <-openPorts:
					fmt.Printf("%s:%d\n", ip, openPort)
				case <-done:
					return
				}
			}
		}()
	}
}
