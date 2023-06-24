package cmd

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/martijnkorbee/gobaboon/tools/boboctl/internal/util"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	butil "github.com/martijnkorbee/gobaboon/pkg/util"
	"github.com/spf13/cobra"
)

var makeNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Make a new baboon web",
	Long:  "Make a new baboon web.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			appName = strings.ToLower(makeNewName)
		)

		// sanitize the application name (convert url to single word)
		if strings.Contains(appName, "/") {
			exploded := strings.SplitAfter(appName, "/")
			appName = exploded[len(exploded)-1]
		}

		// clone the skeleton application
		mustCloneSkeleton(appName)

		// create a ready to go .env file
		mustCreateDotEnv(appName)

		//// update the go.mod file
		//mustUpdateGoModFile(appName)

		// change to new project dir
		err := os.Chdir("./" + appName)
		if err != nil {
			util.PrintFatal("failed to cd to web directory", errors.New("couldn't cd in web directory"))
		}

		// check if not exist create expected folders for a new baboon web
		fnames := []string{
			"cmd/web/web",
			"cmd/web/bin",
			"http/handlers",
			"http/middleware",
			"http/routes",
			"database/models",
			"database/migrations",
			"db-data",
			"templates/views",
			"templates/mail",
			"public/static",
			"logs",
			"tmp",
		}
		for _, fname := range fnames {
			err := butil.CreateDirIfNotExists("./" + fname)
			if err != nil {
				util.PrintWarning("failed to create default dir: " + fname)
			}
		}

		// update the existing .go files with correct name and imports
		color.Yellow("\tUpdating source files...")
		mustUpdateSourceFiles()

		// run go mod tidy in the project directory
		color.Yellow("\tRunning go mod tidy...")

		command := exec.Command("go", "mod", "tidy")
		output, err := command.CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			util.PrintFatal("failed to run go mod tidy", err)
		}

		// TODO: create makefile support for windows
		// update makefile
		mustCreateMakeFile(appName)

		util.PrintSuccess(fmt.Sprint("done creating new web:", appName))

		color.Yellow("\tBuilding: %s", appName)
		command = exec.Command("make", "build_app")
		output, err = command.CombinedOutput()
		if err != nil {
			color.Red(fmt.Sprint(err) + ": " + string(output))
			util.PrintFatal("failed to build web", err)
		}
		util.PrintInfo(fmt.Sprint(string(output)))

		util.PrintSuccess(fmt.Sprintf("Start your web from dir %s and run: /cmd/web/bin/%s", appName, appName))
	},
}

func init() {
	makeNewCmd.Flags().StringVarP(&makeNewName, "name", "n", "", "sets the name of the new web")
	makeNewCmd.MarkFlagRequired("name")
}

func mustCreateDotEnv(appName string) {
	color.Yellow("\tCreating .env file...")

	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		util.PrintFatal("failed to read .env template", err)
	}

	key, err := butil.RandomStringGenerator(32)
	if err != nil {
		util.PrintFatal("failed to make key", err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", key)

	err = os.WriteFile(fmt.Sprintf("./%s/.env", appName), []byte(env), 0644)
	if err != nil {
		util.PrintFatal("failed to write .env file", err)
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
		util.PrintFatal("failed to clone skeleton web", err)
	}

	// remove .git directory
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		util.PrintFatal("failed to remove .git directory", err)
	}
}

func mustUpdateGoModFile(appName string) {
	color.Yellow("\tCreating go.mod file...")

	//_ = os.Remove("./" + appName + "/go.mod")
	//
	//data, err := templateFS.ReadFile("templates/go.mod.txt")
	//if err != nil {
	//	util.PrintFatal("failed to read go mod template", err)
	//}

	data, err := os.ReadFile(fmt.Sprintf("./%s/go.mod", appName))

	mod := string(data)
	mod = strings.ReplaceAll(mod, "github.com/martijnkorbee/gobaboon", makeNewName)

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

		// replace placeholder web name and write new file
		updated := strings.ReplaceAll(string(read), "web", makeNewName)

		err = os.WriteFile(path, []byte(updated), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func mustCreateMakeFile(appName string) {
	color.Yellow("\tUpdating makefile")

	data, err := os.ReadFile("Makefile.unix")
	if err != nil {
		util.PrintFatal("failed to read makefile", err)
	}
	makefile := strings.ReplaceAll(string(data), "${APP_NAME}", appName)
	err = os.WriteFile("Makefile", []byte(makefile), 0644)
	if err != nil {
		util.PrintFatal("failed to write makefile", err)
	}

	err = os.Remove("Makefile.unix")
	if err != nil {
		util.PrintWarning(fmt.Sprint("could not remove Makefile.unix:", err))
	}
}
