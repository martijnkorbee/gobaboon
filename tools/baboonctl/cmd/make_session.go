package cmd

import (
	"fmt"
	ctl "github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var makeSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Make table migrations for baboon sessions",
	Long: `Creates and runs up and down migrations for the baboon server's persistent sessions.
Should be called from the root directory of a baboon app that has a valid .env file.

NOTE: supported databases postgres, mysql/mariadb, sqlite
`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			dbType        = os.Getenv("DATABASE_TYPE")
			migrationName = fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())
			upFilePath    = rootPath + "/database/migrations/" + migrationName + ".up.sql"
			downFilePath  = rootPath + "/database/migrations/" + migrationName + ".down.sql"
		)

		// create up file
		err := ctl.CopyFileFromTemplate(templateFS, "templates/migrations/sessions."+dbType+".up.sql", upFilePath)
		if err != nil {
			ctl.PrintFatal("failed to create up file", err)
		}

		// create down file
		err = ctl.CopyFileFromTemplate(templateFS, "templates/migrations/sessions."+dbType+".down.sql", downFilePath)
		if err != nil {
			ctl.PrintFatal("failed to create down file", err)
		}

		ctl.PrintInfo("created migrations, calling migrate up ...")

		migrateUp()
	},
}
