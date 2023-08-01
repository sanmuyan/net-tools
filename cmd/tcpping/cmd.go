package main

import (
	"flag"
	"fmt"
	"net-tools/pkg/tcpping"
	"time"
)

func main() {
	host := flag.String("h", "localhost", "ping host")
	port := flag.Int("p", 22, "ping port")
	protocol := flag.String("P", "tcp", "ping protocol \nhttp read")
	timeout := flag.Int("T", 1000, "timeout in ms")
	count := flag.Int("c", 4, "ping count")
	interval := flag.Int("i", 1, "ping interval in ms")
	flag.Parse()
	p := tcpping.NewTCPPing(*host, *port, *timeout, *protocol)
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
			fmt.Printf("%s:%d error=%s\n", *host, *port, m)
			errorTotal++
		case t := <-pingTime:
			fmt.Printf("%s:%d time=%dms\n", *host, *port, t)
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
	fmt.Printf("success=%d, error=%d, max=%dms, min=%dms, avg=%dms\n", successTotal, errorTotal, maxTime, minTime, avg)
}
