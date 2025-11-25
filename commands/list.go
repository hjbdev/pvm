package commands

import (
	"fmt"
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

func List() {
	versions, err := common.RetrieveInstalledPHPVersions()
	if err != nil {
		theme.Error(err.Error())
		return
	}

	if len(versions) == 0 {
		theme.Info("No PHP versions installed.")
		return
	}

	var currentVersion string
	if runtime.GOOS == "darwin" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Could not get user home directory: %v", err)
		}
		currentLink := filepath.Join(homeDir, ".pvm", "current")
		linkTarget, err := os.Readlink(currentLink)
		if err == nil {
			currentVersion = filepath.Base(linkTarget)
		}
	} else if runtime.GOOS == "windows" {
		// On Windows, we can check the content of the bat file
		// A more robust solution would be a config file.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Could not get user home directory: %v", err)
		}
		batPath := filepath.Join(homeDir, ".pvm", "bin", "php.bat")
		content, err := os.ReadFile(batPath)
		if err == nil {
			// This is brittle. Assumes path is like "C:\Users\user\.pvm\versions\php-8.1.5..."
			strContent := string(content)
			parts := filepath.SplitList(strContent)
			for _, p := range parts {
				if filepath.Base(p) == "versions" {
					// this logic needs to be better
				}
			}
		}
	}

	theme.Title("Installed PHP versions")
	// print all folders
	for _, version := range versions {
		versionStr := version.StringShort()
		if versionStr == currentVersion {
			color.Green(fmt.Sprintf("  * %s (current)", versionStr))
		} else {
			color.White(fmt.Sprintf("    %s", versionStr))
		}
	}
}
