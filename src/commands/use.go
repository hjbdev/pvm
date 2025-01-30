package commands

import (
	"encoding/json"
	"fmt"
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Use(args []string) {
	threadSafe := true

	if len(args) < 1 {
		theme.Error("You must specify a version or path to use.")
		return
	}

	// If there are two arguments, check for "nts" flag
	if len(args) > 1 && args[1] == "nts" {
		threadSafe = false
	}

	// Get current dir
	currentDir, err := os.Executable()
	if err != nil {
		theme.Error(fmt.Sprintln(err))
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

	// Read versions from versions.json
	versions := readVersionsFromJson(fullDir)
	if versions == nil {
		theme.Error("Error reading versions.json")
		return
	}

	var metaVersions []common.VersionMeta
	// loop over all found installed versions
	for _, version := range *versions {
		safe := true
		if strings.Contains(strings.ToLower(version.Version), "non ") {
			safe = false
		}
		foundVersion := common.ComputeVersion(version.Version, safe, "")
		if threadSafe == foundVersion.ThreadSafe && strings.HasPrefix(foundVersion.String(), args[0]) {
			metaVersions = append(metaVersions, common.VersionMeta{
				Number: foundVersion,
				Folder: version.Path,
			})
		}
	}

	var selectedVersion *common.VersionMeta
	if len(args) == 1 {
		// If a version is provided, prioritize it
		selectedVersion = findVersionByName(metaVersions, args[0], threadSafe)
	} else if len(args) > 1 {
		// If path is provided, find the version by path
		selectedVersion = findVersionByPath(metaVersions, args[0])
	}

	if selectedVersion == nil {
		theme.Error("The specified version or path is not installed.")
		return
	}

	// check if .pvm/bin folder exists
	binPath := filepath.Join(fullDir, "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0755)
	}

	requestedVersion := common.ComputeVersion(args[0], threadSafe, "")
	if requestedVersion.Minor == -1 {
		theme.Warning(fmt.Sprintf("No minor version specified, assumed newest minor version %s.", selectedVersion.Number.String()))
	} else if requestedVersion.Patch == -1 {
		theme.Warning(fmt.Sprintf("No patch version specified, assumed newest patch version %s.", selectedVersion.Number.String()))
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

	// remove the old composer link
	linkPathComposer := filepath.Join(binPath, "composer.phar")
	if _, err := os.Stat(linkPathComposer); err == nil {
		os.Remove(linkPathComposer)
	}

	versionPath := filepath.Join(selectedVersion.Folder, "php.exe")
	versionPathCGI := filepath.Join(selectedVersion.Folder, "php-cgi.exe")
	composerPath := filepath.Join(selectedVersion.Folder, "composer", "composer.phar")

	// create bat script for php
	batCommand := "@echo off \n"
	batCommand = batCommand + "set filepath=\"" + versionPath + "\"\n"
	batCommand = batCommand + "set arguments=%*\n"
	batCommand = batCommand + "%filepath% %arguments%\n"

	err = os.WriteFile(batPath, []byte(batCommand), 0755)

	if err != nil {
		theme.Error(fmt.Sprintln(err))
	}

	// create sh script for php
	shCommand := "#!/bin/bash\n"
	shCommand = shCommand + "filepath=\"" + versionPath + "\"\n"
	shCommand = shCommand + "\"$filepath\" \"$@\""

	err = os.WriteFile(shPath, []byte(shCommand), 0755)

	if err != nil {
		theme.Error(fmt.Sprintln(err))
	}

	// create bat script for php-cgi
	batCommandCGI := "@echo off \n"
	batCommandCGI = batCommandCGI + "set filepath=\"" + versionPathCGI + "\"\n"
	batCommandCGI = batCommandCGI + "set arguments=%*\n"
	batCommandCGI = batCommandCGI + "%filepath% %arguments%\n"

	err = os.WriteFile(batPathCGI, []byte(batCommandCGI), 0755)

	if err != nil {
		theme.Error(fmt.Sprintln(err))
	}

	// create sh script for php-cgi
	shCommandCGI := "#!/bin/bash\n"
	shCommandCGI = shCommandCGI + "filepath=\"" + versionPathCGI + "\"\n"
	shCommandCGI = shCommandCGI + "\"$filepath\" \"$@\""

	err = os.WriteFile(shPathCGI, []byte(shCommandCGI), 0755)

	if err != nil {
		theme.Error(fmt.Sprintln(err))
	}

	// create bat script for composer
	batCommandComposer := "@echo off \n"
	batCommandComposer = batCommandComposer + "set filepath=\"" + versionPath + "\"\n"
	batCommandComposer = batCommandComposer + "set composerpath=\"" + composerPath + "\"\n"
	batCommandComposer = batCommandComposer + "set arguments=%*\n"
	batCommandComposer = batCommandComposer + "%filepath% %composerpath% %arguments%\n"

	err = os.WriteFile(batPathComposer, []byte(batCommandComposer), 0755)

	if err != nil {
		theme.Error(fmt.Sprintln(err))
	}

	// create sh script for php
	shCommandComposer := "#!/bin/bash\n"
	shCommandComposer = shCommandComposer + "filepath=\"" + versionPath + "\"\n"
	shCommandComposer = shCommandComposer + "composerpath=\"" + composerPath + "\"\n"
	shCommandComposer = shCommandComposer + "\"$filepath\" \"$composerpath\" \"$@\""

	err = os.WriteFile(shPathComposer, []byte(shCommandComposer), 0755)

	if err != nil {
		theme.Error(fmt.Sprintln(err))
	}

	var cmd *exec.Cmd
	var output []byte

	composerLinkPath := filepath.Join(binPath, "composer.phar")
	cmd = exec.Command("cmd", "/C", "mklink", "/H", composerLinkPath, composerPath)

	output, err = cmd.Output()
	if err != nil {
		theme.Error(fmt.Sprintln("Error creating composer.phar symlink: %v", err))
		return
	} else {
		theme.Info(string(output))
	}

	// create directory link to ext directory
	extensionDirPath := filepath.Join(selectedVersion.Folder, "ext")
	extensionLinkPath := filepath.Join(binPath, "ext")

	// delete the old link first if it exists
	if _, err := os.Stat(extensionLinkPath); err == nil {
		cmd := exec.Command("cmd", "/C", "rmdir", extensionLinkPath)
		_, err := cmd.Output()
		if err != nil {
			theme.Error(fmt.Sprintf("Error deleting ext directory directory link: %v", err))
			return
		}
	}

	// create directory link - uses cmd since using os.Symlink did require extra permissions
	cmd = exec.Command("cmd", "/C", "mklink", "/J", extensionLinkPath, extensionDirPath)

	output, err = cmd.Output()
	if err != nil {
		theme.Error(fmt.Sprintln("Error creating ext directory symlink: %v", err))
		return
	} else {
		theme.Info(string(output))
	}
	// end of ext directory link creation

	theme.Success(fmt.Sprintf("Using PHP %s", selectedVersion.Number))
}

// readVersionsFromJson reads and returns the versions from versions.json
func readVersionsFromJson(fullDir string) *[]common.VersionJson {
	// Path to the versions.json file
	jsonPath := filepath.Join(fullDir, "versions", "versions.json")

	// Read the existing versions from versions.json
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil
	}

	// Expecting the JSON content to be a string
	var existingVersions []common.VersionJson
	if err := json.Unmarshal(data, &existingVersions); err != nil {
		return nil
	}

	return &existingVersions
}

// findVersionByName finds and returns a version by its name
func findVersionByName(versions []common.VersionMeta, versionName string, threadSafe bool) *common.VersionMeta {
	for _, version := range versions {
		// Match by version name and thread-safe flag
		if strings.HasPrefix(version.Number.String(), versionName) && version.Number.ThreadSafe == threadSafe {
			return &version
		}
	}
	return nil
}

// findVersionByPath finds and returns a version by its path
func findVersionByPath(versions []common.VersionMeta, path string) *common.VersionMeta {
	for _, version := range versions {
		// Match the path in the version folder
		versionFolderPath := filepath.Join(version.Folder, "php.exe")
		if strings.Contains(versionFolderPath, path) {
			return &version
		}
	}
	return nil
}
