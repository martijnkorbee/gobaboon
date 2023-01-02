package cmd

import (
	"embed"
	"os"

	"github.com/joho/godotenv"
	"github.com/martijnkorbee/gobaboon/cmd/cli/internal/pkg/util"
	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"
)

var (
	//go:embed templates
	templateFS embed.FS

	rootpath string
	dotenv   bool
)

var rootCmd = &cobra.Command{
	Use:   "bobo",
	Short: "Bobo is a quality of life improvement for building web apps and services.",
	Long:  `Bobo is a convenient tool to bootstrap your web apps and services. Write your first routes within 10 min.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.PrintWarning("no command specified, use bobo help")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	// set rootpath for bobo
	if path, err := os.Getwd(); err != nil {
		util.PrintFatal("failed to get working directory", err)
	} else {
		rootpath = path
	}

	// check load .env
	if util.FileExists(rootpath + "/.env") {
		err := godotenv.Load(rootpath + "/.env")
		if err != nil {
			util.PrintFatal("failed to load .env", err)
		}
		dotenv = true
	}
}
