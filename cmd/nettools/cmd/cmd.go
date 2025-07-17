package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net-tools/pkg/nettestc"
	"net-tools/pkg/nettests"
	"net-tools/pkg/portscan"
	"net-tools/pkg/speedtestc"
	"net-tools/pkg/speedtests"
	"net-tools/pkg/tcpping"
)

var rootCtx context.Context

var rootCmd = &cobra.Command{
	Use:   "net-tools",
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
}

var tcpPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "TCP ping",
	Run: func(cmd *cobra.Command, args []string) {
		tcpping.Run(rootCtx, args)
	},
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

var netTestsCmd = &cobra.Command{
	Use:   "nts",
	Short: "Network test server",
	Run: func(cmd *cobra.Command, args []string) {
		nettests.Run(rootCtx)
	},
}

var netTestcCmd = &cobra.Command{
	Use:   "ntc",
	Short: "Network test client",
	Run: func(cmd *cobra.Command, args []string) {
		nettestc.Run(rootCtx)
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

	portScanCmd.Flags().IntP("max-thread", "t", 1, "max thread")
	portScanCmd.Flags().IntP("timeout", "T", 200, "timeout (ms)")

	tcpPingCmd.Flags().StringP("protocol", "P", "tcp", "protocol (tcp|http)")
	tcpPingCmd.Flags().IntP("timeout", "T", 1000, "timeout (ms)")
	tcpPingCmd.Flags().IntP("count", "C", 4, "count")
	tcpPingCmd.Flags().IntP("interval", "i", 1000, "interval (ms)")
	tcpPingCmd.Flags().Bool("tls", false, "with tls")

	speedTestsCmd.Flags().StringP("server-bind", "s", ":8080", "server bind addr")

	speedTestcCmd.Flags().StringP("protocol", "P", "tcp", "test protocol (tcp|quic)")
	speedTestcCmd.Flags().StringP("server-addr", "s", "localhost:8080", "server addr")
	speedTestcCmd.Flags().IntP("test-time", "t", 10, "test time (s)")
	speedTestcCmd.Flags().StringP("test-mode", "m", "download", "test mode (download|upload)")
	speedTestcCmd.Flags().IntP("max-thread", "T", 1, "test max thread")

	netTestsCmd.Flags().StringP("server-bind", "s", ":8080", "server bind addr")
	netTestsCmd.Flags().StringP("protocol", "P", "tcp", "test protocol (tcp|udp|http|ws)")
	netTestsCmd.Flags().Uint32P("timeout", "t", 1000*5, "client timeout (ms)")

	netTestcCmd.Flags().StringP("server-addr", "s", "localhost:8080", "server addr")
	netTestcCmd.Flags().StringP("protocol", "P", "tcp", "test protocol (tcp|udp|http|ws)")
	netTestcCmd.Flags().Uint32P("timeout", "T", 1000*5, "server timeout (ms)")
	netTestcCmd.Flags().Uint32P("interval", "i", 1000, "test interval (ms)")
	netTestcCmd.Flags().Uint32P("max-thread", "t", 1, "test max thread")

	rootCmd.AddCommand(portScanCmd, tcpPingCmd, speedTestsCmd, speedTestcCmd, netTestsCmd, netTestcCmd)
}

func Execute(ctx context.Context) {
	rootCtx = ctx
	if err := rootCmd.Execute(); err != nil {
		logrus.Tracef("cmd execute error: %v", err)
	}
}
