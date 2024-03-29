package cmd

import (
	ctl "github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"github.com/spf13/cobra"
)

var migrateForceCmd = &cobra.Command{
	Use:   "force",
	Short: "Resets migration version",
	Long: `Resets migration version in case of a dirty version, if version is specified forces specified version.
Without a specified version after you can migrate up any step you want.
NOTE: Force does not run any actual migrations. Without a specified version migrate down won't work after forcing.`,
	Run: func(cmd *cobra.Command, args []string) {
		if migrationForceVersion > 0 {
			migrateForceVersion(migrationForceVersion)
		} else {
			migrateForce()
		}
	},
}

func init() {
	migrateForceCmd.Flags().IntVarP(&migrationForceVersion, "version", "v", 0, "forces a specific migration version")
}

func migrateForce() {
	m := mustConnectMigrator()
	defer m.Close()

	if err := m.Force(-1); err != nil {
		ctl.PrintError("error forcing migration", err)
	}

	ctl.PrintSuccess("force migration successfully")
}

func migrateForceVersion(v int) {
	m := mustConnectMigrator()
	defer m.Close()

	if err := m.Force(v); err != nil {
		ctl.PrintError("error forcing migration version", err)
	}

	ctl.PrintResult("force migration successfully version", v)
}
