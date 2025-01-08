package commands

import (
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"log"
	"slices"

	"github.com/fatih/color"
)

func ListRemote() {
	versions, err := common.RetrievePHPVersions()
	if err != nil {
		log.Fatalln(err)
	}

	common.SortVersions(versions)

	installedVersions, _ := common.RetrieveInstalledPHPVersions()

	theme.Title("PHP versions available")
	for _, version := range versions {
		idx := slices.IndexFunc(installedVersions, func(v common.Version) bool { return v.Same(version) })
		found := " "
		if idx != -1 {
			found = "*"
		}
		color.White(found + "   v" + version.StringShort())
	}
}
