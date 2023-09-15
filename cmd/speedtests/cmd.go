package main

import (
	"flag"
	"net-tools/pkg/speedtests"
)

func main() {
	serverBind := flag.String("s", "0.0.0.0", "SpeedTest Server bind")
	serverPort := flag.Int("p", 8080, "SpeedTest Server port")
	protocol := flag.String("P", "tcp", "SpeedTest protocol (tcp|udp)")
	flag.Parse()
	speedtests.Start(speedtests.NewServer(*serverBind, *serverPort, *protocol))
}
