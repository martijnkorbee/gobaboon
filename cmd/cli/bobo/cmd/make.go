package cmd

import (
	"github.com/spf13/cobra"
)

var (
	makeModelMigrate bool
	makeModelName    string
	makeNewName      string
)

var makeCmd = &cobra.Command{
	Use:       "make",
	Short:     "Make all kinds of things",
	Long:      "Make all kinds of things with the Bobo cli.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"key", "session", "auth", "model"},
}

func init() {
	rootCmd.AddCommand(makeCmd)

	makeCmd.AddCommand(
		makeKeyCmd,
		makeSessionCmd,
		makeAuthCmd,
		makeModelCmd,
		makeNewCmd,
	)
}
