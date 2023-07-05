package cmd

import (
	"github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"github.com/spf13/cobra"
)

var rpcMaintenanceCmd = &cobra.Command{
	Use:   "maintenance",
	Short: "Set the server in maintenance mode",
	Long:  "Set the server in maintenance mode.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			resp string
		)

		// dial rpc server
		c, err := dialRPC()
		if err != nil {
			util.PrintError("failed to make rpc call", err)
			return
		}

		err = c.Call("RPCServer.SetMaintenanceMode", rpcMaintenanceOn, &resp)
		if err != nil {
			util.PrintError("failed to execute call", err)
			return
		}

		util.PrintSuccess(resp)
	},
}

func init() {
	rpcMaintenanceCmd.Flags().BoolVarP(&rpcMaintenanceOn, "mode", "m", false, "set maintenance mode")
	rpcMaintenanceCmd.MarkFlagRequired("mode")
}
