package cmd

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	butil "github.com/martijnkorbee/gobaboon/pkg/util"
	"github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var makeNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Make a new app",
	Long:  "Make a new app.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			appName = strings.ToLower(makeNewName)
		)

		// sanitize the application name (convert url to single word)
		if strings.Contains(appName, "/") {
			exploded := strings.SplitAfter(appName, "/")
			appName = exploded[len(exploded)-1]
		}

		// clone baboon project
		mustCloneSkeleton(appName)

		// create the .config.properties file
		mustCreateConfig(appName)

		// update the go.mod file
		mustUpdateGoModFile(appName)

		// change to new project dir
		err := os.Chdir("./" + appName)
		if err != nil {
			util.PrintFatal("failed to cd to app directory", errors.New("failed to cd in app directory"))
		}

		// update the existing .go files with correct name and imports
		color.Yellow("\tUpdating source files...")
		mustUpdateSourceFiles()

		// remove existing LICENSE file
		color.Yellow("\tDeleting LICENSE...")
		if err = os.Remove("LICENSE"); err != nil {
			util.PrintWarning(err.Error())
		}

		// run go mod tidy in the project directory
		color.Yellow("\tRunning go mod tidy...")

		command := exec.Command("go", "mod", "tidy")
		output, err := command.CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			util.PrintFatal("failed to run go mod tidy", err)
		}

		util.PrintSuccess(fmt.Sprint("done creating new application with name:", appName))

		color.Yellow("\tBuilding: %s", appName)
		command = exec.Command("make", "app_build")
		output, err = command.CombinedOutput()
		if err != nil {
			color.Red(fmt.Sprint(err) + ": " + string(output))
			util.PrintFatal("failed to build app", err)
		}
		util.PrintInfo(fmt.Sprint(string(output)))

		util.PrintSuccess(fmt.Sprintf("Go to the new %s directory and type: make app_start", appName))
	},
}

func init() {
	makeNewCmd.Flags().StringVarP(&makeNewName, "name", "n", "", "sets the name of the new app")
	makeNewCmd.MarkFlagRequired("name")
}

func mustCreateConfig(appName string) {
	color.Yellow("\tCreating .config.properties file...")

	data, err := templateFS.ReadFile("templates/config.example")
	if err != nil {
		util.PrintFatal("failed to read .config.example template", err)
	}

	key, err := butil.RandomStringGenerator(32)
	if err != nil {
		util.PrintFatal("failed to make key", err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", key)

	err = os.WriteFile(fmt.Sprintf("./%s/app/.config.properties", appName), []byte(env), 0644)
	if err != nil {
		util.PrintFatal("failed to write .config.properties file", err)
	}
}

func mustCloneSkeleton(appName string) {
	color.Yellow("\tCloning git repository...")

	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/martijnkorbee/gobaboon.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		util.PrintFatal("failed to clone skeleton app", err)
	}

	// remove .git directory
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		util.PrintFatal("failed to remove .git directory", err)
	}
}

func mustUpdateGoModFile(appName string) {
	color.Yellow("\tCreating go.mod file...")

	data, err := os.ReadFile(fmt.Sprintf("./%s/go.mod", appName))

	mod := string(data)
	mod = strings.ReplaceAll(mod, "github.com/martijnkorbee/gobaboon", appName)

	err = os.WriteFile(fmt.Sprintf("./%s/go.mod", appName), []byte(mod), 0644)
	if err != nil {
		util.PrintFatal("failed write go mod file", err)
	}
}

func mustUpdateSourceFiles() {
	// walk entire project folder, including subfolders
	err := filepath.Walk(".", updateSoureFile)
	if err != nil {
		util.PrintFatal("failed to update source file", err)
	}
}

func updateSoureFile(path string, fi os.FileInfo, err error) error {
	// check for an error before doing anything else
	if err != nil {
		return err
	}

	// check if file is directory
	if fi.IsDir() {
		return nil
	}

	// only check go files
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}

	if matched {
		// read file
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// replace placeholder app name and write new file
		updated := strings.ReplaceAll(string(read), "github.com/martijnkorbee/gobaboon", makeNewName)

		err = os.WriteFile(path, []byte(updated), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
