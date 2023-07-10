package cmd

import (
	"errors"
	"fmt"
	ctl "github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"os"
	"strings"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
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
	},
}

func init() {
	makeModelCmd.Flags().StringVarP(&makeModelName, "name", "n", "", "name of the new model")
	makeModelCmd.Flags().StringVarP(&dbType, "db-type", "t", "", "specify your database type")

	makeModelCmd.MarkFlagRequired("name")
	makeModelCmd.MarkFlagRequired("db-type")
}

func mustMakeModel(modelName, tableName string) {
	var (
		fileName = rootPath + "/internal/database/models/" + strings.ToLower(modelName) + ".go"
	)

	// check if model doesn't exist already
	if ctl.FileExists(fileName) {
		ctl.PrintFatal("failed to make model", errors.New("model already exists"))
		return
	}

	// read model template
	data, err := templateFS.ReadFile("templates/models/model.go.txt")
	if err != nil {
		ctl.PrintFatal("failed to read model template", err)
	}

	// update template with new model name
	model := string(data)
	model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

	err = os.WriteFile(fileName, []byte(model), 0644)
	if err != nil {
		ctl.PrintFatal("failed to write model in models directory", err)
	}

	// feedback
	ctl.PrintResult("created model:", fileName)
}

func mustMakeModelMigrations(tableName string) {
	var (
		migrationName = fmt.Sprintf("%d_create_%s_table", time.Now().UnixMicro(), tableName)
		upSource      = "templates/migrations/model_table." + dbType + ".up.sql.txt"
		downSource    = "templates/migrations/model_table." + dbType + ".down.sql.txt"
		upTarget      = rootPath + "/internal/database/migrations/" + migrationName + ".up.sql"
		downTarget    = rootPath + "/internal/database/migrations/" + migrationName + ".down.sql"
	)

	// create up file
	data, err := templateFS.ReadFile(upSource)
	if err != nil {
		ctl.PrintFatal("failed to read migration template", err)
	}
	// update migration template with new model name
	tmpl := string(data)
	tmpl = strings.ReplaceAll(tmpl, "$TABLENAME$", tableName)
	// write migration
	if err := os.WriteFile(upTarget, []byte(tmpl), 0644); err != nil {
		ctl.PrintFatal("failed to write migration in migrations directory", err)
	}

	// create down file
	data, err = templateFS.ReadFile(downSource)
	if err != nil {
		ctl.PrintFatal("failed to read migration template", err)
	}
	// update migration template with new model name
	tmpl = string(data)
	tmpl = strings.ReplaceAll(tmpl, "$TABLENAME$", tableName)
	// write migration
	if err := os.WriteFile(downTarget, []byte(tmpl), 0644); err != nil {
		ctl.PrintFatal("failed to write migration in migrations directory", err)
	}

	// feedback
	ctl.PrintResult("created migrations for model", tableName)
}
