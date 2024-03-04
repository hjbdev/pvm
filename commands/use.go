package commands

import (
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Use(args []string) {
	threadSafe := true

	if len(args) < 1 {
		theme.Error("You must specify a version to use.")
		return
	}

	if len(args) > 1 {
		if args[1] == "nts" {
			threadSafe = false
		}
	}

	// get users home dir
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln(err)
	}

	// check if .pvm folder exists
	if _, err := os.Stat(filepath.Join(homeDir, ".pvm")); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// check if .pvm/versions folder exists
	if _, err := os.Stat(filepath.Join(homeDir, ".pvm", "versions")); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// check if .pvm/bin folder exists
	if _, err := os.Stat(filepath.Join(homeDir, ".pvm", "bin")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(homeDir, ".pvm", "bin"), 0755)
	}

	// get all folders in .pvm/versions
	versions, err := os.ReadDir(filepath.Join(homeDir, ".pvm", "versions"))
	if err != nil {
		log.Fatalln(err)
	}

	// check if version exists
	selectedVersion := ""
	for _, version := range versions {
		versionNumbers := common.GetVersion(version.Name())
		if versionNumbers.Major+"."+versionNumbers.Minor+"."+versionNumbers.Patch == args[0] {
			if threadSafe && !strings.Contains(version.Name(), "nts") {
				selectedVersion = version.Name()
			} else if !threadSafe && strings.Contains(version.Name(), "nts") {
				selectedVersion = version.Name()
			}
		}
	}

	if selectedVersion == "" {
		theme.Error("The specified version is not installed.")
		return
	}

	// remove old symlink
	linkPath := filepath.Join(homeDir, ".pvm", "bin", "php")
	if _, err := os.Stat(linkPath); err == nil {
		os.Remove(linkPath)
	}

	// create symlink
	versionPath := filepath.Join(homeDir, ".pvm", "versions", selectedVersion, "php.exe")

	cmd := exec.Command("cmd", "/c", "mklink", "/H", linkPath, versionPath)
	err = cmd.Run()
	if err != nil {
		theme.Error("Failed to create symlink")
		return
	}

	var threadSafeString string
	if threadSafe {
		threadSafeString = "non-thread safe"
	} else {
		threadSafeString = "thread safe"
	}

	theme.Success("Using PHP " + args[0] + " " + threadSafeString)
}
