package cmd

import (
	"github.com/martijnkorbee/gobaboon/tools/boboctl/internal/util"
	"github.com/spf13/cobra"
)

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Runs all non run down migrations",
	Long:  "Runs all non run down migrations.",
	Run: func(cmd *cobra.Command, args []string) {
		if migrationSteps > 0 {
			migrateDownSteps(-migrationSteps)
		} else {
			migrateDown()
		}
	},
}

func init() {
	migrateDownCmd.Flags().IntVarP(&migrationSteps, "steps", "n", 0, "migration steps")
}

func migrateDown() {
	m := mustConnectMigrator()
	defer m.Close()

	if err := m.Down(); err != nil {
		util.PrintError("error running down migration", err)
		return
	}

	util.PrintSuccess("all down migrations ran successfully")
}

func migrateDownSteps(n int) {
	m := mustConnectMigrator()
	defer m.Close()

	if err := m.Steps(n); err != nil {
		util.PrintError("error running down migration", err)
		return
	}

	util.PrintSuccess("down migrations ran successfully")
}
