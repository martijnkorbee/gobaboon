package cmd

import (
	"net/rpc"
	"os"

	"github.com/spf13/cobra"
)

var (
	rpcMaintenanceOn bool
)

var rpcCmd = &cobra.Command{
	Use:       "rpc",
	Short:     "Used to make rpc calls to the baboon app",
	Long:      "Used to make rpc calls to the baboon app.",
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"maintenance"},
}

func init() {
	rootCmd.AddCommand(rpcCmd)

	rpcCmd.AddCommand(rpcMaintenanceCmd)
}

func dialRPC() (*rpc.Client, error) {
	var (
		rpcHost = os.Getenv("RPC_HOST")
		rpcPort = os.Getenv("RPC_PORT")
	)

	// rpc host default is 0.0.0.0
	if rpcHost == "" {
		rpcHost = "0.0.0.0"
	}

	// rpc default port is 4004
	if rpcPort == "" {
		rpcPort = "4004"
	}

	if c, err := rpc.Dial("tcp", "rpcHost:"+rpcPort); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}
