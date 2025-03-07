package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"net-tools/cmd/nettools/cmd"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigs
		logrus.Debug("process is shutting down...")
		cancel()
	}()
	cmd.Execute(ctx)
}
