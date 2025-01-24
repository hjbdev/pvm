package commands

import (
	"fmt"
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Uninstall(args []string) {
	threadSafe := true

	if len(args) < 1 {
		theme.Error("You must specify a version to uninstall.")
		return
	}

	if len(args) > 1 {
		if args[1] == "nts" {
			threadSafe = false
		}
	}

	// get current dir
	currentDir, err := os.Executable()

	if err != nil {
		log.Fatalln(err)
		return
	}

	fullDir := filepath.Dir(currentDir)

	// check if .pvm folder exists
	if _, err := os.Stat(filepath.Join(fullDir)); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// check if .pvm/versions folder exists
	if _, err := os.Stat(filepath.Join(fullDir, "versions")); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// check if .pvm/bin folder exists
	binPath := filepath.Join(fullDir, "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0755)
	}

	// get all folders in .pvm/versions
	versions, err := os.ReadDir(filepath.Join(fullDir, "versions"))
	if err != nil {
		log.Fatalln(err)
	}

	var selectedVersion *common.VersionMeta
	// loop over all found installed versions
	for i, version := range versions {
		safe := true
		if strings.Contains(version.Name(), "nts") || strings.Contains(version.Name(), "NTS") {
			safe = false
		}
		foundVersion := common.ComputeVersion(version.Name(), safe, "")
		if threadSafe == foundVersion.ThreadSafe && strings.HasPrefix(foundVersion.String(), args[0]) {
			selectedVersion = &common.VersionMeta{
				Number: foundVersion,
				Folder: versions[i],
			}
		}
	}

	if selectedVersion == nil {
		theme.Error("The specified version is not installed.")
		return
	}

	// remove old php bat script
	batPath := filepath.Join(binPath, "php.bat")
	if _, err := os.Stat(batPath); err == nil {
		os.Remove(batPath)
	}

	// remove the old php sh script
	shPath := filepath.Join(binPath, "php")
	if _, err := os.Stat(shPath); err == nil {
		os.Remove(shPath)
	}

	// remove old php-cgi bat script
	batPathCGI := filepath.Join(binPath, "php-cgi.bat")
	if _, err := os.Stat(batPathCGI); err == nil {
		os.Remove(batPathCGI)
	}

	// remove old php-cgi sh script
	shPathCGI := filepath.Join(binPath, "php-cgi")
	if _, err := os.Stat(shPathCGI); err == nil {
		os.Remove(shPathCGI)
	}

	// remove old composer bat script
	batPathComposer := filepath.Join(binPath, "composer.bat")
	if _, err := os.Stat(batPathComposer); err == nil {
		os.Remove(batPathComposer)
	}

	// remove the old composer sh script
	shPathComposer := filepath.Join(binPath, "composer")
	if _, err := os.Stat(shPathComposer); err == nil {
		os.Remove(shPathComposer)
	}

	// create directory link to ext directory
	extensionLinkPath := filepath.Join(binPath, "ext")

	// delete the old link first if it exists
	if _, err := os.Stat(extensionLinkPath); err == nil {
		cmd := exec.Command("cmd", "/C", "rmdir", extensionLinkPath)
		_, err := cmd.Output()
		if err != nil {
			log.Fatalln("Error deleting ext directory directory link:", err)
			return
		}
	}

	versionFolderPath := filepath.Join(fullDir, "versions", selectedVersion.Folder.Name())
	if _, err := os.Stat(versionFolderPath); err == nil {
		os.RemoveAll(versionFolderPath)
	}

	theme.Success(fmt.Sprintf("Uninstalling PHP %s", selectedVersion.Number))
}
