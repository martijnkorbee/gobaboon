package cmd

import (
	"log"
	"os"

	"github.com/martijnkorbee/gobaboon/cmd/cli/internal/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "creates readme docs",
	Run: func(cmd *cobra.Command, args []string) {

		rootpath, err := os.Getwd()
		if err != nil {
			util.PrintError("failed to get rootpath", err)
		}

		err = doc.GenMarkdownTree(rootCmd, rootpath+"/cmd/cli/docs")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
