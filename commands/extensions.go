package commands

import (
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Extensions(args []string) {
	if len(args) < 2 {
		theme.Error("You must specify an action and an extension.")
		theme.Info("Usage: pvm extensions <enable|disable> <extension>")
		return
	}

	// determine which version is currently selected
	currentVersion := common.GetCurrentVersionFolder()

	if currentVersion == "" {
		theme.Error("You do not have an active PHP version.")
		theme.Info("Select a PHP version with `pvm use <version>` first.")
	}

	command := args[0]
	ext := args[1]

	if command != "enable" && command != "disable" {
		theme.Error("Invalid action. Must be 'enable' or 'disable'.")
		return
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	// check if the version exists
	if _, err := os.Stat(homeDir + "/.pvm/versions/" + currentVersion); os.IsNotExist(err) {
		theme.Error("The specified version does not exist.")
		return
	}

	extensions := strings.Split(ext, ",")

	for _, extension := range extensions {
		handleExtension(extension, command, homeDir, currentVersion)
	}
}

func handleExtension(ext string, command string, homeDir string, currentVersion string) {
	ini := common.ReadPhpIni(homeDir + "/.pvm/versions/" + currentVersion + "/php.ini")
	splitIni := regexp.MustCompile(`\r?\n`).Split(ini, -1)
	extensionStatus, lineNumber := common.GetExtensionStatus(ini, ext)

	if extensionStatus == common.ExtensionEnabled {
		if command == "enable" {
			theme.Success("Extension " + ext + " is already enabled.")
		} else {
			disabledLine := ";" + splitIni[lineNumber]
			splitIni[lineNumber] = disabledLine
			newIni := strings.Join(splitIni, "\n")
			os.WriteFile(filepath.Join(homeDir, ".pvm", "versions", currentVersion, "php.ini"), []byte(newIni), 0644)
			theme.Success("Extension " + ext + " enabled.")
		}
	} else if extensionStatus == common.ExtensionDisabled {
		if command == "enable" {
			enabledLine := strings.Replace(splitIni[lineNumber], ";", "", 1)
			splitIni[lineNumber] = enabledLine
			newIni := strings.Join(splitIni, "\n")
			os.WriteFile(filepath.Join(homeDir, ".pvm", "versions", currentVersion, "php.ini"), []byte(newIni), 0644)
			theme.Success("Extension " + ext + " enabled.")
		} else {
			theme.Success("Extension " + ext + " is already disabled.")
		}
	} else {
		theme.Error("Extension " + ext + " not found in php.ini")
	}
}
