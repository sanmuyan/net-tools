package main

import (
	"flag"
	"net-tools/pkg/speedtestc"
	"os"
	"strconv"
	"strings"
)

func main() {
	testTime := flag.Int("t", 10, "SpeedTest time (seconds)")
	mode := flag.String("m", "download", "SpeedTest mode (download|upload)")
	protocol := flag.String("P", "tcp", "SpeedTest protocol (tcp|udp)")
	maxThread := flag.Int("T", 1, "SpeedTest Max thread")
	server := flag.String("s", "localhost:8080", "SpeedTest server")
	flag.Parse()

	if flag.NArg() == 1 {
		*server = flag.Arg(0)
	}
	addr := strings.Split(*server, ":")
	if len(addr) != 2 {
		println("Invalid server address: ", os.Args[1])
		println("Example: speedtestc 192.168.1.1:8080")
		os.Exit(1)
	}
	host := addr[0]
	port, _ := strconv.Atoi(addr[1])
	speedtestc.Start(speedtestc.NewClient(host, port, *mode, *testTime, *protocol, *maxThread))
}
