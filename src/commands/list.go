package commands

import (
	"hjbdev/pvm/theme"
	"log"
	"os"
	"github.com/fatih/color"
	"path/filepath"
)

func List() {
	// get current dir
	currentDir, err := os.Executable()

	if err != nil {
		log.Fatalln(err)
	}
	
	fullDir := filepath.Dir(currentDir)

	if err != nil {
		log.Fatalln(err)
	}

	// check if .pvm folder exists
	if _, err := os.Stat(fullDir); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// check if .pvm/versions folder exists
	if _, err := os.Stat(fullDir + "/versions"); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// get all folders in .pvm/versions
	versions, err := os.ReadDir(fullDir + "/versions")
	if err != nil {
		log.Fatalln(err)
	}

	theme.Title("Installed PHP versions")

	// print all folders
	for _, version := range versions {
		color.White("    " + version.Name())
	}
}
