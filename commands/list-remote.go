package commands

import (
	"fmt"
	"os"

	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"slices"

	"github.com/fatih/color"
)

var retrievePHPVersions = common.RetrievePHPVersions
var retrieveInstalledPHPVersions = common.RetrieveInstalledPHPVersions

func ListRemote() error {
	versions, err := retrievePHPVersions()
	if err != nil {
		return err
	}

	common.SortVersions(versions)

	installedVersions, _ := retrieveInstalledPHPVersions()

	// Only attempt to read the current-version metadata when HOME is set
	// (tests set HOME when they want to provide a controlled environment).
	var currentVersion string
	var currentVersionNumber common.Version
	var currentVersionErr error
	if os.Getenv("HOME") != "" {
		currentVersion = common.GetCurrentVersionFolder()
		currentVersionNumber, currentVersionErr = common.ParseVersion(currentVersion, common.IsThreadSafeName(currentVersion), "")
	} else {
		currentVersionErr = fmt.Errorf("HOME not set")
	}

	theme.Title("PHP versions available")
	for _, version := range versions {
		label := version.StringShort()

		idx := slices.IndexFunc(installedVersions, func(v common.Version) bool { return v.Same(version) })
		if idx != -1 {
			if currentVersionErr == nil && version.Same(currentVersionNumber) {
				label += " " + color.GreenString("[current]")
			} else {
				label += " " + color.CyanString("[installed]")
			}
		}

		color.White("    " + label)
	}

	return nil
}
