package cmd

import (
	ctl "github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
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
		ctl.PrintError("error running down migration", err)
		return
	}

	ctl.PrintSuccess("all down migrations ran successfully")
}

func migrateDownSteps(n int) {
	m := mustConnectMigrator()
	defer m.Close()

	if err := m.Steps(n); err != nil {
		ctl.PrintError("error running down migration", err)
		return
	}

	ctl.PrintSuccess("down migrations ran successfully")
}
