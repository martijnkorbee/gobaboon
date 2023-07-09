package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"github.com/spf13/cobra"
)

var makeModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Make a new model",
	Long: `Creates a new model in the models directory and respective migrations in migrations directory,
if the migrate flag is passed also runs up migrations.

NOTE: supported databases postgres, mysql/mariadb, sqlite
`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			tableName, modelName string
			// create a new pluralize instance
			plur = pluralize.NewClient()
		)

		// sanitize modelname if necessary
		if plur.IsPlural(makeModelName) {
			modelName = plur.Singular(makeModelName)
			tableName = strings.ToLower(makeModelName)
		} else {
			modelName = makeModelName
			tableName = strings.ToLower(plur.Plural(makeModelName))
		}

		// create the new model
		mustMakeModel(modelName, tableName)

		// make migrations
		mustMakeModelMigrations(tableName)

		// run migration
		if makeModelMigrate {
			migrateUp()
		}
	},
}

func init() {
	makeModelCmd.Flags().StringVarP(&makeModelName, "name", "n", "", "name of the new model")
	makeModelCmd.Flags().BoolVarP(&makeModelMigrate, "migrate", "m", false, "set to auto run migrations")
	// required flags
	makeModelCmd.MarkFlagRequired("name")
}

func mustMakeModel(modelName, tableName string) {
	var (
		fileName = rootpath + "/database/models/" + strings.ToLower(modelName) + ".go"
	)

	// check if model doesn't exists
	if util.FileExists(fileName) {
		util.PrintFatal("failed to make model", errors.New("model already exists"))
		return
	}

	// read model go text file
	data, err := templateFS.ReadFile("templates/models/model.go.txt")
	if err != nil {
		util.PrintFatal("failed to read model template", err)
	}

	// update model go text file with new model name
	model := string(data)
	model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

	err = os.WriteFile(fileName, []byte(model), 0644)
	if err != nil {
		util.PrintFatal("failed to write model in models directory", err)
	}

	// print model feedback
	util.PrintResult("created model", fileName)
}

func mustMakeModelMigrations(tableName string) {
	// check .env
	if !dotenv {
		util.PrintError("failed to make migrations for model", errors.New("no .env file in current directory"))
		return
	}

	var (
		dbType        = os.Getenv("DATABASE_TYPE")
		migrationName = fmt.Sprintf("%d_create_%s_table", time.Now().UnixMicro(), tableName)
		upFilePath    = rootpath + "/database/migrations/" + migrationName + ".up.sql"
		downFilePath  = rootpath + "/database/migrations/" + migrationName + ".down.sql"
	)

	// create up file
	data, err := templateFS.ReadFile("templates/migrations/model_table." + dbType + ".up.sql.txt")
	if err != nil {
		util.PrintFatal("failed to read migration template", err)
	}
	// update migration template with new model name
	tmpl := string(data)
	tmpl = strings.ReplaceAll(tmpl, "$TABLENAME$", tableName)
	// write migration
	err = os.WriteFile(upFilePath, []byte(tmpl), 0644)
	if err != nil {
		util.PrintFatal("failed to write migration in migrations directory", err)
	}

	// create down file
	data, err = templateFS.ReadFile("templates/migrations/model_table." + dbType + ".down.sql.txt")
	if err != nil {
		util.PrintFatal("failed to read migration template", err)
	}
	// update migration template with new model name
	tmpl = string(data)
	tmpl = strings.ReplaceAll(tmpl, "$TABLENAME$", tableName)
	// write migration
	err = os.WriteFile(downFilePath, []byte(tmpl), 0644)
	if err != nil {
		util.PrintFatal("failed to write migration in migrations directory", err)
	}

	// feedback
	util.PrintResult("created migrations for model", tableName)
}
