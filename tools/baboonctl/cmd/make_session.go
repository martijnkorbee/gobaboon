package cmd

import (
	"fmt"
	ctl "github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"time"

	"github.com/spf13/cobra"
)

var makeSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Make table migrations for baboon sessions",
	Long: `Creates up and down migrations for the server's persistent sessions.

SUPPORTED DATABASES: [postgres, mysql, mariadb, sqlite]
`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			migrationName = fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())
			upSource      = "templates/migrations/sessions." + dbType + ".up.sql"
			downSource    = "templates/migrations/sessions." + dbType + ".down.sql"
			upTarget      = rootPath + "/internal/database/migrations/" + migrationName + ".up.sql"
			downTarget    = rootPath + "/internal/database/migrations/" + migrationName + ".down.sql"
		)

		// create up file
		err := ctl.CopyFileFromTemplate(templateFS, upSource, upTarget)
		if err != nil {
			ctl.PrintFatal("failed to create up file", err)
		}

		// create down file
		err = ctl.CopyFileFromTemplate(templateFS, downSource, downTarget)
		if err != nil {
			ctl.PrintFatal("failed to create down file", err)
		}

		ctl.PrintResult("created migrations for", "persistent sessions")
	},
}

func init() {
	makeSessionCmd.Flags().StringVarP(&dbType, "db-type", "t", "", "specify your database type")
	makeSessionCmd.MarkFlagRequired("db-type")
}
