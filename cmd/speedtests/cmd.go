package main

import (
	"flag"
	"net-tools/pkg/speedtests"
)

func main() {
	serverBind := flag.String("s", "0.0.0.0", "server bind")
	serverPort := flag.Int("p", 8080, "server port")
	protocol := flag.String("P", "tcp", "test protocol tcp or udp")
	flag.Parse()
	speedtests.Start(speedtests.NewServer(*serverBind, *serverPort, *protocol))
}
