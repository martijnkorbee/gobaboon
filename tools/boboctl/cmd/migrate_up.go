package cmd

import (
	"github.com/martijnkorbee/gobaboon/tools/boboctl/internal/util"
	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Runs all non run up migrations",
	Long:  "Runs all non run up migrations.",
	Run: func(cmd *cobra.Command, args []string) {
		if migrationSteps > 0 {
			migrateUpSteps(migrationSteps)
		} else {
			migrateUp()
		}
	},
}

func init() {
	migrateUpCmd.Flags().IntVarP(&migrationSteps, "steps", "n", 0, "migration steps")
}

func migrateUp() {
	m := mustConnectMigrator()
	defer m.Close()

	if err := m.Up(); err != nil {
		util.PrintError("error running up migration", err)
		return
	}

	util.PrintSuccess("all up migrations ran successfully")
}

func migrateUpSteps(n int) {
	m := mustConnectMigrator()
	defer m.Close()

	if err := m.Steps(n); err != nil {
		util.PrintError("error running up migration", err)
		return
	}

	util.PrintSuccess("up migrations ran successfully")
}
