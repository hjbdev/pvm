package commands

import (
	"fmt"
	"hjbdev/pvm/theme"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Uninstall(args []string) {
	if len(args) < 2 {
		theme.Error("You must specify a version to uninstall.")
		return
	}

	version := args[1]

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not get user home directory: %v", err)
	}
	versionsPath := filepath.Join(homeDir, ".pvm", "versions")

	if runtime.GOOS == "darwin" {
		// Check if the version is installed

		uninstallForMac(version, versionsPath)
	} else if runtime.GOOS == "windows" {
		uninstallForWindows(version, versionsPath)
	} else {
		theme.Error(fmt.Sprintf("Unsupported operating system: %s", runtime.GOOS))
	}
}

func uninstallForMac(version string, versionsPath string) {
	versionPath := filepath.Join(versionsPath, version)
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		theme.Error(fmt.Sprintf("PHP version %s is not installed.", version))
		return
	}
	theme.Info(fmt.Sprintf("Uninstalling PHP %s for macOS using Homebrew...", version))

	phpPackage := fmt.Sprintf("php@%s", version)
	cmd := exec.Command("brew", "uninstall", phpPackage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		theme.Warning(fmt.Sprintf("Could not uninstall %s using Homebrew. It might have been installed manually or uninstalled already.", phpPackage))
		log.Printf("Brew uninstall error: %v\n", err)
	}

	theme.Info(fmt.Sprintf("Removing version directory: %s", versionPath))
	if err := os.RemoveAll(versionPath); err != nil {
		theme.Error(fmt.Sprintf("Failed to remove directory for PHP %s.", version))
		log.Fatalf("Error: %v", err)
	}

	theme.Success(fmt.Sprintf("Successfully uninstalled PHP %s.", version))
}

// isCompleteVersion checks if the version is complete (3 parts)
func isCompleteVersion(version string) bool {
	parts := strings.Split(version, ".")
	return len(parts) == 3
}

func uninstallForWindows(version string, versionsPath string) {
	theme.Info(fmt.Sprintf("Uninstalling PHP %s for Windows...", version))
	// Validate version format and prompt for complete version if needed
	isCompleteVersion := isCompleteVersion(version)
	if !isCompleteVersion {
		log.Fatalf("Version '%s' is incomplete. Please provide the complete version (e.g., %s.1, %s.2, etc.)", version, version, version)
	}
	// For windows, the version folder name is like php-8.4.14-Win32-vs17-x64
	// We need to find the directory that contains the version number
	files, err := os.ReadDir(versionsPath)
	if err != nil {
		log.Fatalf("Could not read versions directory: %v", err)
	}

	var dirToRemove string
	for _, f := range files {
		if f.IsDir() {
			// Check if the directory name contains the version number
			// Windows folder names are like: php-8.4.14-Win32-vs17-x64
			// We want to match when user provides "8.4.14"
			dirName := f.Name()
			// Look for patterns like php-8.4.14- or just 8.4.14 in the directory name
			if strings.Contains(dirName, version) {
				// Additional check to ensure it's a proper version match
				// This prevents partial matches like "8.4" matching "8.4.14"
				versionPattern := fmt.Sprintf("php-%s-", version)
				if strings.HasPrefix(dirName, versionPattern) || dirName == version {
					dirToRemove = filepath.Join(versionsPath, dirName)
					break
				}
			}
		}
	}

	if dirToRemove == "" {
		theme.Error(fmt.Sprintf("Could not find installation directory for version %s", version))
		return
	}

	theme.Info(fmt.Sprintf("Removing directory: %s", dirToRemove))
	if err := os.RemoveAll(dirToRemove); err != nil {
		theme.Error(fmt.Sprintf("Failed to remove directory for PHP %s.", version))
		log.Fatalf("Error: %v", err)
	}

	theme.Success(fmt.Sprintf("Successfully uninstalled PHP %s.", version))
}
