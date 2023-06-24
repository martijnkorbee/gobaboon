package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/martijnkorbee/gobaboon/cmd/boboctl/internal/pkg/util"
	"github.com/spf13/cobra"
)

var (
	migrationSteps        int
	migrationForceVersion int
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Create or run migrations",
	Long:  "Create or run migrations.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.AddCommand(
		migrateUpCmd,
		migrateDownCmd,
		migrateForceCmd,
	)
}

func mustConnectMigrator() *migrate.Migrate {
	var (
		source = "file://" + rootpath + "/database/migrations"
		dsn    = mustBuildDSN()
	)

	m, err := migrate.New(source, dsn)
	if err != nil {
		util.PrintFatal("could not connect to database", err)
	}

	return m
}

func mustBuildDSN() string {
	var (
		dbType = os.Getenv("DATABASE_TYPE")
		dsn    string
	)

	switch dbType {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASS"),
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
		)
	case "mysql", "mariadb":
		dsn = fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?tls=%s",
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASS"),
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
		)
	case "sqlite":
		dsn = fmt.Sprintf("sqlite3://%s/database/data/sqlite/%s.db",
			rootpath,
			os.Getenv("DATABASE_NAME"),
		)
	default:
		util.PrintFatal("failed to build connection string", errors.New("unsupported database type"))
	}

	return dsn
}
