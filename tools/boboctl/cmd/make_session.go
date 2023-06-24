package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/martijnkorbee/gobaboon/cmd/boboctl/internal/pkg/util"
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
		// check .env
		if !dotenv {
			util.PrintError("failed to make sessions", errors.New("no .env file in current directory"))
			return
		}

		var (
			dbType        = os.Getenv("DATABASE_TYPE")
			migrationName = fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())
			upFilePath    = rootpath + "/database/migrations/" + migrationName + ".up.sql"
			downFilePath  = rootpath + "/database/migrations/" + migrationName + ".down.sql"
		)

		// create up file
		err := util.CopyFileFromTemplate(templateFS, "templates/migrations/sessions."+dbType+".up.sql", upFilePath)
		if err != nil {
			util.PrintFatal("failed to create up file", err)
		}

		// create down file
		err = util.CopyFileFromTemplate(templateFS, "templates/migrations/sessions."+dbType+".down.sql", downFilePath)
		if err != nil {
			util.PrintFatal("failed to create down file", err)
		}

		util.PrintInfo("created migrations, calling migrate up ...")

		migrateUp()
	},
}
