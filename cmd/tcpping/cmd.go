package main

import (
	"flag"
	"fmt"
	"github.com/sanmuyan/xpkg/xnet"
	"net"
	"net-tools/pkg/tcpping"
	"os"
	"strconv"
	"time"
)

func main() {
	protocol := flag.String("P", "tcp", "Ping protocol (tcp|http|read)")
	timeout := flag.Int("T", 1000, "Ping Timeout (ms)")
	count := flag.Int("c", 4, "Ping count")
	interval := flag.Int("i", 1, "Ping interval (ms)")
	flag.Parse()

	if len(os.Args) < 3 {
		println("Example: tcpping 192.168.1.1 22")
		os.Exit(1)
	}
	host := os.Args[1]
	_port := os.Args[2]
	if !xnet.IsIP(host) {
		_, err := net.LookupHost(host)
		if err != nil {
			println("ping: Name or service not known", host)
			os.Exit(1)
		}
	}
	if !xnet.IsPort(_port) {
		println("ping: Invalid port: ", _port)
		os.Exit(1)
	}
	port, _ := strconv.Atoi(_port)
	p := tcpping.NewTCPPing(host, port, *timeout, *protocol)
	errorMessage := make(chan string)
	pingTime := make(chan int)
	go func() {
		for i := 0; i < *count; i++ {
			p.PING(errorMessage, pingTime)
			time.Sleep(time.Duration(*interval) * time.Millisecond)
		}
	}()
	var totalTime int
	var successTotal int
	var errorTotal int
	var maxTime int
	var minTime int
	for i := 0; i < *count; i++ {
		select {
		case m := <-errorMessage:
			fmt.Printf("Reply from %s:%d error=%s\n", host, port, m)
			errorTotal++
		case t := <-pingTime:
			fmt.Printf("Reply from %s:%d time=%dms\n", host, port, t)
			totalTime += t
			successTotal++
			if t > maxTime {
				maxTime = t
			}
			if t < minTime || minTime == 0 {
				minTime = t
			}
		}
	}
	var avg int
	if successTotal > 0 {
		avg = totalTime / successTotal
	}
	fmt.Printf("Success=%d, Error=%d, Max=%dms, Min=%dms, Avg=%dms\n", successTotal, errorTotal, maxTime, minTime, avg)
}
