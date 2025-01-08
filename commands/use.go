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
	binPath := filepath.Join(homeDir, ".pvm", "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0755)
	}

	// get all folders in .pvm/versions
	versions, err := os.ReadDir(filepath.Join(homeDir, ".pvm", "versions"))
	if err != nil {
		log.Fatalln(err)
	}

	var selectedVersion *versionMeta
	// loop over all found installed versions
	for i, version := range versions {
		safe := true
		if strings.Contains(version.Name(), "nts") || strings.Contains(version.Name(), "NTS") {
			safe = false
		}
		foundVersion := common.ComputeVersion(version.Name(), safe, "")
		if threadSafe == foundVersion.ThreadSafe && strings.HasPrefix(foundVersion.String(), args[0]) {
			selectedVersion = &versionMeta{
				number: foundVersion,
				folder: versions[i],
			}
		}
	}

	if selectedVersion == nil {
		theme.Error("The specified version is not installed.")
		return
	}

	requestedVersion := common.ComputeVersion(args[0], threadSafe, "")
	if requestedVersion.Minor == -1 {
		theme.Warning(fmt.Sprintf("No minor version specified, assumed newest minor version %s.", selectedVersion.number.String()))
	} else if requestedVersion.Patch == -1 {
		theme.Warning(fmt.Sprintf("No patch version specified, assumed newest patch version %s.", selectedVersion.number.String()))
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

	versionFolderPath := filepath.Join(homeDir, ".pvm", "versions", selectedVersion.folder.Name())
	versionPath := filepath.Join(versionFolderPath, "php.exe")
	versionPathCGI := filepath.Join(versionFolderPath, "php-cgi.exe")

	// create bat script for php
	batCommand := "@echo off \n"
	batCommand = batCommand + "set filepath=\"" + versionPath + "\"\n"
	batCommand = batCommand + "set arguments=%*\n"
	batCommand = batCommand + "%filepath% %arguments%\n"

	err = os.WriteFile(batPath, []byte(batCommand), 0755)

	if err != nil {
		log.Fatalln(err)
	}

	// create sh script for php
	shCommand := "#!/bin/bash\n"
	shCommand = shCommand + "filepath=\"" + versionPath + "\"\n"
	shCommand = shCommand + "\"$filepath\" \"$@\""

	err = os.WriteFile(shPath, []byte(shCommand), 0755)

	if err != nil {
		log.Fatalln(err)
	}

	// create bat script for php-cgi
	batCommandCGI := "@echo off \n"
	batCommandCGI = batCommandCGI + "set filepath=\"" + versionPathCGI + "\"\n"
	batCommandCGI = batCommandCGI + "set arguments=%*\n"
	batCommandCGI = batCommandCGI + "%filepath% %arguments%\n"

	err = os.WriteFile(batPathCGI, []byte(batCommandCGI), 0755)

	if err != nil {
		log.Fatalln(err)
	}

	// create sh script for php-cgi
	shCommandCGI := "#!/bin/bash\n"
	shCommandCGI = shCommandCGI + "filepath=\"" + versionPathCGI + "\"\n"
	shCommandCGI = shCommandCGI + "\"$filepath\" \"$@\""

	err = os.WriteFile(shPathCGI, []byte(shCommandCGI), 0755)

	if err != nil {
		log.Fatalln(err)
	}

	// create directory link to ext directory
	extensionDirPath := filepath.Join(versionFolderPath, "ext")
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

	// create directory link - uses cmd since using os.Symlink did require extra permissions
	cmd := exec.Command("cmd", "/C", "mklink", "/J", extensionLinkPath, extensionDirPath)

	output, err := cmd.Output()
	if err != nil {
		log.Fatalln("Error creating ext directory symlink:", err)
		return
	} else {
		theme.Info(string(output))
	}
	// end of ext directory link creation

	theme.Success(fmt.Sprintf("Using PHP %s", selectedVersion.number))
}

type versionMeta struct {
	number common.Version
	folder os.DirEntry
}
