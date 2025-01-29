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

// Uninstall will remove the specified PHP version and its entry from versions.json.
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

	// Get current directory
	currentDir, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
		return
	}

	fullDir := filepath.Dir(currentDir)

	// Check if .pvm folder exists
	if _, err := os.Stat(filepath.Join(fullDir)); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// Check if .pvm/versions folder exists
	if _, err := os.Stat(filepath.Join(fullDir, "versions")); os.IsNotExist(err) {
		theme.Error("No PHP versions installed")
		return
	}

	// Check if .pvm/bin folder exists
	binPath := filepath.Join(fullDir, "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0755)
	}

	// Get all folders in .pvm/versions
	versions, err := os.ReadDir(filepath.Join(fullDir, "versions"))
	if err != nil {
		log.Fatalln(err)
	}

	var selectedVersion *common.VersionMeta
	// Loop over all found installed versions
	for i, version := range versions {
		safe := true
		if strings.Contains(version.Name(), "nts") || strings.Contains(version.Name(), "NTS") {
			safe = false
		}
		foundVersion := common.ComputeVersion(version.Name(), safe, "")
		if threadSafe == foundVersion.ThreadSafe && strings.HasPrefix(foundVersion.String(), args[0]) {
			selectedVersion = &common.VersionMeta{
				Number: foundVersion,
				Folder: versions[i].Name(),
			}
		}
	}

	if selectedVersion == nil {
		theme.Error("The specified version is not installed.")
		return
	}

	// Remove php bat and sh scripts
	scripts := []string{"php.bat", "php", "php-cgi.bat", "php-cgi", "composer.bat", "composer"}

	for _, script := range scripts {
		scriptPath := filepath.Join(binPath, script)
		if _, err := os.Stat(scriptPath); err == nil {
			if err := os.Remove(scriptPath); err != nil {
				theme.Error(fmt.Sprintf("Error removing script %s: %v", script, err))
				return
			}
		}
	}

	// Create directory link to ext directory
	extensionLinkPath := filepath.Join(binPath, "ext")

	// Delete the old link first if it exists
	if _, err := os.Stat(extensionLinkPath); err == nil {
		cmd := exec.Command("cmd", "/C", "rmdir", extensionLinkPath)
		_, err := cmd.Output()
		if err != nil {
			theme.Error(fmt.Sprintf("Error deleting ext directory link: %v", err))
			return
		}
	}

	// Remove the version folder from versions directory
	versionFolderPath := filepath.Join(fullDir, "versions", selectedVersion.Folder)
	if _, err := os.Stat(versionFolderPath); err == nil {
		os.RemoveAll(versionFolderPath)
	}

	// Remove the version entry from versions.json
	if err := common.RemoveFromVersionJson(versionFolderPath); err != nil {
		theme.Error(fmt.Sprintf("Error removing version from versions.json: %v", err))
	}

	theme.Success(fmt.Sprintf("Finished uninstalling PHP %s", selectedVersion.Number))
}
