package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net-tools/pkg/benchtestc"
	"net-tools/pkg/benchtests"
	"net-tools/pkg/portscan"
	"net-tools/pkg/speedtestc"
	"net-tools/pkg/speedtests"
	"net-tools/pkg/tcpping"
)

var rootCtx context.Context

var rootCmd = &cobra.Command{
	Use:   "nts",
	Short: "Network tools",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := initConfig(cmd)
		if err != nil {
			logrus.Fatalf("init config error: %v", err)
		}
	},
}

var portScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Port scan",
	Run: func(cmd *cobra.Command, args []string) {
		portscan.Run(rootCtx, args)
	},
	Example: "192.168.1.1 1-65535",
}

var tcpPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "TCP ping",
	Run: func(cmd *cobra.Command, args []string) {
		tcpping.Run(rootCtx, args)
	},
	Example: "192.168.1.1 80",
}

var speedTestsCmd = &cobra.Command{
	Use:   "sts",
	Short: "Speed test server",
	Run: func(cmd *cobra.Command, args []string) {
		speedtests.Run(rootCtx)
	},
}

var speedTestcCmd = &cobra.Command{
	Use:   "stc",
	Short: "Speed test client",
	Run: func(cmd *cobra.Command, args []string) {
		speedtestc.Run(rootCtx)
	},
}

var benchTestsCmd = &cobra.Command{
	Use:   "bts",
	Short: "Benchmark test server",
	Run: func(cmd *cobra.Command, args []string) {
		benchtests.Run(rootCtx)
	},
}

var benchTestcCmd = &cobra.Command{
	Use:   "btc",
	Short: "Benchmark test client",
	Run: func(cmd *cobra.Command, args []string) {
		benchtestc.Run(rootCtx)
	},
}

var configFile string

const (
	logLevel = 4
)

func init() {
	// 初始化命令行参数
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().IntP("log-level", "l", logLevel, "log level")
	rootCmd.PersistentFlags().BoolP("pprof-server", "", false, "enable pprof server")

	portScanCmd.Flags().IntP("max-thread", "T", 1, "max thread")
	portScanCmd.Flags().IntP("timeout", "t", 200, "timeout (ms)")

	tcpPingCmd.Flags().StringP("protocol", "P", "tcp", "protocol (tcp|http)")
	tcpPingCmd.Flags().IntP("timeout", "t", 1000, "timeout (ms)")
	tcpPingCmd.Flags().IntP("count", "C", 4, "count")
	tcpPingCmd.Flags().IntP("interval", "i", 1000, "interval (ms)")
	tcpPingCmd.Flags().Bool("tls", false, "with tls")

	speedTestsCmd.Flags().StringP("server-bind", "s", ":8080", "server bind addr")

	speedTestcCmd.Flags().StringP("protocol", "P", "tcp", "test protocol (tcp|quic)")
	speedTestcCmd.Flags().StringP("server-addr", "s", "localhost:8080", "server addr")
	speedTestcCmd.Flags().IntP("test-time", "t", 10, "test time (s)")
	speedTestcCmd.Flags().StringP("test-mode", "m", "download", "test mode (download|upload)")
	speedTestcCmd.Flags().IntP("max-thread", "T", 1, "test max thread")

	benchTestsCmd.Flags().StringP("server-bind", "s", ":8080", "server bind addr")
	benchTestsCmd.Flags().StringP("protocol", "P", "tcp", "test protocol (tcp|udp|http|https|ws)")
	benchTestsCmd.Flags().IntP("timeout", "t", 1000*60, "read message timeout (ms)")

	benchTestcCmd.Flags().StringP("server-addr", "s", "localhost:8080", "server addr")
	benchTestcCmd.Flags().StringP("protocol", "P", "tcp", "test protocol (tcp|udp|http|https|ws)")
	benchTestcCmd.Flags().IntP("timeout", "t", 1000, "read message timeout (ms)")
	benchTestcCmd.Flags().IntP("interval", "i", 1000, "send message interval (ms)")
	benchTestcCmd.Flags().IntP("max-thread", "T", 1, "test max thread")
	benchTestcCmd.Flags().IntP("max-messages", "m", 0, "max send messages")

	rootCmd.AddCommand(portScanCmd, tcpPingCmd, speedTestsCmd, speedTestcCmd, benchTestsCmd, benchTestcCmd)
}

func Execute(ctx context.Context) {
	rootCtx = ctx
	if err := rootCmd.Execute(); err != nil {
		logrus.Tracef("cmd execute error: %v", err)
	}
}
