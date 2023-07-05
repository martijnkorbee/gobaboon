package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"github.com/spf13/cobra"
)

var makeAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Make table migrations for baboon auth",
	Long: `Creates and runs up and down migrations for the baboon auth tables, and adds user and token models in data directory.
Should be called from the root directory of a baboon web that has a valid .env file.

NOTE: supported databases postgres, mysql/mariadb, sqlite
`,
	Run: func(cmd *cobra.Command, args []string) {
		// check .env
		if !dotenv {
			util.PrintError("failed to make auth tables", errors.New("no .env file in current directory"))
			return
		}

		var (
			dbType        = os.Getenv("DATABASE_TYPE")
			migrationName = fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
			upFilePath    = rootpath + "/database/migrations/" + migrationName + ".up.sql"
			downFilePath  = rootpath + "/database/migrations/" + migrationName + ".down.sql"
		)

		// create up file
		err := util.CopyFileFromTemplate(templateFS, "templates/migrations/auth_tables."+dbType+".up.sql", upFilePath)
		if err != nil {
			util.PrintFatal("failed to create up file", err)
		}

		// create down file
		err = util.CopyFileFromTemplate(templateFS, "templates/migrations/auth_tables."+dbType+".down.sql", downFilePath)
		if err != nil {
			util.PrintFatal("failed to create down file", err)
		}

		util.PrintInfo("created migrations, calling migrate up ...")

		// run migrations
		migrateUp()

		// create user and token models in data directory
		err = util.CopyFileFromTemplate(templateFS, "templates/data/token.go.txt", rootpath+"/database/models/token.go")
		if err != nil {
			util.PrintError("failed to create token model", err)
		}

		err = util.CopyFileFromTemplate(templateFS, "templates/data/user.go.txt", rootpath+"/database/models/user.go")
		if err != nil {
			util.PrintError("failed to create user model", err)
		}

		// create auth middleware in middleware directory
		err = util.CopyFileFromTemplate(templateFS, "templates/middleware/auth-token.go.txt", rootpath+"/http/middleware/auth-token.go")
		if err != nil {
			util.PrintError("failed to create auth token middleware", err)
		}

		err = util.CopyFileFromTemplate(templateFS, "templates/middleware/auth-user.go.txt", rootpath+"/http/middleware/auth-user.go")
		if err != nil {
			util.PrintError("failed to create auth user middleware", err)
		}
	},
}
