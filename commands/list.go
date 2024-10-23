package commands

import (
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"log"
	"os"

	"github.com/fatih/color"
)

func List() {
	// get users home dir
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln(err)
	}

	// check if .pvm folder exists
	if _, err := os.Stat(homeDir + "/.pvm"); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// check if .pvm/versions folder exists
	if _, err := os.Stat(homeDir + "/.pvm/versions"); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// get all folders in .pvm/versions
	versions, err := os.ReadDir(homeDir + "/.pvm/versions")
	if err != nil {
		log.Fatalln(err)
	}

	theme.Title("Installed PHP versions")

	currentVersion := common.GetCurrentVersionFolder()

	// print all folders
	for _, version := range versions {
		if version.Name() == currentVersion {
			color.White("    " + version.Name() + " (current)")
		} else {
			color.White("    " + version.Name())
		}
	}
}
