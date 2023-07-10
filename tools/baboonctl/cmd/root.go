package cmd

import (
	"embed"
	ctl "github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"os"

	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"
)

var (
	//go:embed templates
	templateFS embed.FS

	rootPath string
	dbType   string
	dbName   string
)

var rootCmd = &cobra.Command{
	Use:   "baboonctl",
	Short: "Baboonctl is a quality of life improvement for building app apps and services.",
	Long:  `Baboonctl is a convenient tool to bootstrap your app apps and services. Write your first routes in minutes.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctl.PrintWarning("no command specified, use bobo help")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// set rootPath for baboonctl
	if path, err := os.Getwd(); err != nil {
		ctl.PrintFatal("failed to get working directory", err)
	} else {
		rootPath = path
	}
}
