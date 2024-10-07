package commands

import (
	"fmt"
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

	// transform to easily sortable slice
	var availableVersions []versionMeta
	for i, version := range versions {
		availableVersions = append(availableVersions, versionMeta{
			number: common.GetVersion(version.Name()),
			folder: versions[i],
		})
	}

	// check if version exists
	var selectedVersion *versionMeta
	for _, version := range availableVersions {
		if fmt.Sprintf("%d.%d.%d", version.number.Major, version.number.Minor, version.number.Patch) == args[0] {
			if threadSafe && !strings.Contains(version.folder.Name(), "nts") {
				selectedVersion = &versionMeta{
					number: version.number,
					folder: version.folder,
				}
			} else if !threadSafe && strings.Contains(version.folder.Name(), "nts") {
				selectedVersion = &versionMeta{
					number: version.number,
					folder: version.folder,
				}
			}
		}
	}

	// if patch version is not specified, use the newest matching major.minor
	if selectedVersion == nil {
		// Sort by newest patch first
		availableVersions = sortVersions(availableVersions)

		for _, version := range availableVersions {
			if fmt.Sprintf("%d.%d", version.number.Major, version.number.Minor) == args[0] {
				if threadSafe && !strings.Contains(version.folder.Name(), "nts") {
					selectedVersion = &versionMeta{
						number: version.number,
						folder: version.folder,
					}
				} else if !threadSafe && strings.Contains(version.folder.Name(), "nts") {
					selectedVersion = &versionMeta{
						number: version.number,
						folder: version.folder,
					}
				}
				break
			}
		}

		if selectedVersion == nil {
			theme.Error("The specified version is not installed.")
			return
		} else {
			theme.Warning(fmt.Sprintf("No patch version specified, assumed newest patch version %s.", selectedVersion.number.String()))
		}
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

	var threadSafeString string
	if threadSafe {
		threadSafeString = "thread safe"
	} else {
		threadSafeString = "non-thread safe"
	}

	theme.Success("Using PHP " + selectedVersion.number.String() + " " + threadSafeString)
}

func sortVersions(in []versionMeta) []versionMeta {
	sort.Slice(in, func(i, j int) bool {
		if in[i].number.Major != in[j].number.Major {
			return in[i].number.Major > in[j].number.Major
		}
		if in[i].number.Minor != in[j].number.Minor {
			return in[i].number.Minor > in[j].number.Minor
		}
		return in[i].number.Patch > in[j].number.Patch
	})

	return in
}

type versionMeta struct {
	number common.Version
	folder os.DirEntry
}
