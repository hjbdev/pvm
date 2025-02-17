package commands

import (
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"github.com/fatih/color"
)

func List() {
	versions, err := common.RetrieveInstalledPHPVersions()
	if err != nil {
		theme.Error(err.Error())
	}

	theme.Title("Installed PHP versions")
	// print all folders
	for _, version := range versions {
		color.White("    " + version.StringShort())
	}
}
