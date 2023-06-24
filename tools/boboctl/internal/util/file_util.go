package util

import (
	"embed"
	"os"
)

// CopyFileFromTemplate takes in a template FS full template and target path,
// checks if the file does not exist and copies the file.
func CopyFileFromTemplate(fs embed.FS, template string, target string) error {
	// skip if target file exists
	if FileExists(target) {
		PrintWarning("skipped file, already exists: " + target)
		return nil
	}

	tmpl, err := fs.ReadFile(template)
	if err != nil {
		PrintFatal("could not read template", err)
	}

	err = os.WriteFile(target, []byte(tmpl), 0644)
	if err != nil {
		PrintFatal("failed to write file", err)
	}

	PrintResult("created", target)

	return nil
}

// FileExists takes a full file path and checks if the file exists and returns true or false.
func FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
