package main

import (
	"flag"
	"fmt"
	"github.com/sanmuyan/xpkg/xnet"
	"net"
	"net-tools/pkg/tcpping"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	protocol := flag.String("P", "tcp", "Ping protocol (tcp|http|https|read)")
	timeout := flag.Int("T", 1000, "Ping Timeout (ms)")
	count := flag.Int("c", 4, "Ping count")
	interval := flag.Int("i", 1000, "Ping interval (ms)")
	pingHost := flag.String("h", "127.0.0.1:80", "Ping host")
	flag.Parse()

	var host string
	var port string
	if flag.NArg() == 1 {
		*pingHost = flag.Arg(0)
	}
	hostAndPort := strings.Split(*pingHost, ":")
	if len(hostAndPort) != 2 {
		println("Example: tcpping 127.0.0.1:80")
		os.Exit(1)
	}
	host = hostAndPort[0]
	port = hostAndPort[1]
	if !xnet.IsIP(host) {
		_, err := net.LookupHost(host)
		if err != nil {
			println("ping: Name or service not known", host)
			println("Example: tcpping 127.0.0.1:80")
			os.Exit(1)
		}
	}
	if !xnet.IsPort(port) {
		println("ping: Invalid port: ", port)
		println("Example: tcpping 127.0.0.1:80")
		os.Exit(1)
	}
	portInt, _ := strconv.Atoi(port)
	p := tcpping.NewTCPPing(host, portInt, int64(*timeout), *protocol)
	errorMessage := make(chan string)
	pingTime := make(chan float32)
	go func() {
		for i := 0; i < *count; i++ {
			p.PING(errorMessage, pingTime)
			time.Sleep(time.Duration(*interval) * time.Millisecond)
		}
	}()
	var totalTime float32
	var successTotal int
	var errorTotal int
	var maxTime float32
	var minTime float32
	for i := 0; i < *count; i++ {
		select {
		case m := <-errorMessage:
			fmt.Printf("Reply from %s:%d error=%s\n", host, portInt, m)
			errorTotal++
		case t := <-pingTime:
			t = t / 1000
			fmt.Printf("Reply from %s:%d time=%.3fms\n", host, portInt, t)
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
	var avg float32
	if successTotal > 0 {
		avg = totalTime / float32(successTotal)
	}
	fmt.Printf("Success=%d, Error=%d, Max=%.3fms, Min=%.3fms, Avg=%.3fms\n", successTotal, errorTotal, maxTime, minTime, avg)
}
