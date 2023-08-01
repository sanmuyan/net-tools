package main

import (
	"flag"
	"net-tools/pkg/speedtestc"
)

func main() {
	serverHost := flag.String("s", "localhost", "server host")
	serverPort := flag.Int("p", 8080, "server port")
	testTime := flag.Int("t", 10, "test time in seconds")
	mode := flag.String("m", "download", "download or upload")
	protocol := flag.String("P", "tcp", "test protocol tcp or udp")
	maxThread := flag.Int("T", 1, "max thread")
	flag.Parse()
	speedtestc.Start(speedtestc.NewClient(*serverHost, *serverPort, *mode, *testTime, *protocol, *maxThread))
}
