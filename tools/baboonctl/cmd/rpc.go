package cmd

import (
	"errors"
	"net/rpc"
	"os"

	"github.com/spf13/cobra"
)

var (
	rpcMaintenanceOn bool
)

var rpcCmd = &cobra.Command{
	Use:       "rpc",
	Short:     "Used to make rpc calls to the baboon web",
	Long:      "Used to make rpc calls to the baboon web.",
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"maintenance"},
}

func init() {
	rootCmd.AddCommand(rpcCmd)

	rpcCmd.AddCommand(rpcMaintenanceCmd)
}

func dialRPC() (*rpc.Client, error) {
	// check .env
	if !dotenv {
		return nil, errors.New("no .env file in current directory")
	}

	var (
		rpcport = os.Getenv("RPC_PORT")
	)

	// rpc default port is 4004
	if rpcport == "" {
		rpcport = "4004"
	}

	if c, err := rpc.Dial("tcp", "127.0.0.1:"+rpcport); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}
