package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Bobo",
	Long:  `All software has versions. This is Bobo's.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.HiWhite("Bobo version: " + version)
	},
}
