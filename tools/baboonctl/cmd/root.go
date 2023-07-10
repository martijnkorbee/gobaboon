package cmd

import (
	"embed"
	"github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"os"

	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"
)

var (
	//go:embed templates
	templateFS embed.FS

<<<<<<< HEAD
	rootPath string
	dbType   string
	dbName   string
=======
	rootpath string
	dotenv   bool
	dbtype   string
>>>>>>> main
)

var rootCmd = &cobra.Command{
	Use:   "baboonctl",
	Short: "Baboonctl is a quality of life improvement for building app apps and services.",
	Long:  `Baboonctl is a convenient tool to bootstrap your app apps and services. Write your first routes in minutes.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.PrintWarning("no command specified, use bobo help")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

<<<<<<< HEAD
	rootCmd.PersistentFlags().StringVarP(&dbType, "db-type", "D", "", "specify your database type")
	rootCmd.PersistentFlags().StringVarP(&dbName, "db-name", "N", "", "specify your database name")
}

func initConfig() {
	// set rootPath for baboonctl
	if path, err := os.Getwd(); err != nil {
		util.PrintFatal("failed to get working directory", err)
	} else {
		rootPath = path
=======
	rootCmd.PersistentFlags().StringVarP(&dbtype, "db-type", "D", "", "specify your database type")
}

func initConfig() {
	// set rootpath for baboonctl
	if path, err := os.Getwd(); err != nil {
		util.PrintFatal("failed to get working directory", err)
	} else {
		rootpath = path
>>>>>>> main
	}
}
