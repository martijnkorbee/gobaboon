package cmd

import (
	"fmt"
	"time"

	"github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"github.com/spf13/cobra"
)

var makeAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Make table migrations for baboon auth",
	Long: `Creates up and down migrations for the auth tables, and adds user and token models in models directory.
Should be called from the root directory of a your application.
SUPPORTED DATABASES: postgres, mysql/mariadb, sqlite
`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			migrationName = fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())

			// migrations
			upSource   = "templates/migrations/auth_tables." + dbType + ".up.sql"
			downSource = "templates/migrations/auth_tables." + dbType + ".down.sql"
			upTarget   = rootPath + "/internal/database/migrations/" + migrationName + ".up.sql"
			downTarget = rootPath + "/internal/database/migrations/" + migrationName + ".down.sql"

			// models
			tokenSource = "templates/models/token.go.txt"
			usersSource = "templates/models/user.go.txt"
			tokenTarget = rootPath + "/internal/database/models/token.go"
			usersTarget = rootPath + "/internal/database/models/user.go"

			// middleware
			authTokenSource = "templates/middleware/auth-token.go.txt"
			authUsersSource = "templates/middleware/auth-user.go.txt"
			authTokenTarget = rootPath + "/internal/http/middleware/auth-token.go"
			authUsersTarget = rootPath + "/internal/http/middleware/auth-user.go"
		)

		// create database migrations
		if err := util.CopyFileFromTemplate(templateFS, upSource, upTarget); err != nil {
			util.PrintFatal("failed to create up file", err)
		}
		if err := util.CopyFileFromTemplate(templateFS, downSource, downTarget); err != nil {
			util.PrintFatal("failed to create down file", err)
		}

		// create database models
		if err := util.CopyFileFromTemplate(templateFS, tokenSource, tokenTarget); err != nil {
			util.PrintError("failed to create token model", err)
		}
		if err := util.CopyFileFromTemplate(templateFS, usersSource, usersTarget); err != nil {
			util.PrintError("failed to create user model", err)
		}

		// create auth middleware
		if err := util.CopyFileFromTemplate(templateFS, authTokenSource, authTokenTarget); err != nil {
			util.PrintError("failed to create auth token middleware", err)
		}

		if err := util.CopyFileFromTemplate(templateFS, authUsersSource, authUsersTarget); err != nil {
			util.PrintError("failed to create auth user middleware", err)
		}
	},
}
