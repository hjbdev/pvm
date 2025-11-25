package commands

import (
	"fmt"
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Use(args []string) {
	if len(args) < 1 {
		theme.Error("You must specify a version to use.")
		return
	}

	if runtime.GOOS == "darwin" {
		useForMac(args)
	} else if runtime.GOOS == "windows" {
		useForWindows(args)
	} else {
		theme.Error(fmt.Sprintf("Unsupported operating system: %s", runtime.GOOS))
	}
}

func useForMac(args []string) {
	version := args[0]
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not get user's home directory: %v", err)
	}

	pvmPath := filepath.Join(homeDir, ".pvm")
	versionsPath := filepath.Join(pvmPath, "versions")
	versionPath := filepath.Join(versionsPath, version)

	// Check if the version is installed
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		theme.Error(fmt.Sprintf("PHP %s is not installed. Run 'pvm install %s' to install it.", version, version))
		return
	}

	// Symlink logic
	currentLinkPath := filepath.Join(pvmPath, "current")

	// Remove existing symlink
	if _, err := os.Lstat(currentLinkPath); err == nil {
		if err := os.Remove(currentLinkPath); err != nil {
			log.Fatalf("Failed to remove existing 'current' symlink: %v", err)
		}
	}

	// Create new symlink
	if err := os.Symlink(versionPath, currentLinkPath); err != nil {
		log.Fatalf("Failed to create 'current' symlink for %s: %v", version, err)
	}

	// Check shell config for PATH
	shell_config := os.Getenv("SHELL")
	if shell_config == "" {
		shell_config = "bash" // default
	}

	var shellrc string
	if strings.Contains(shell_config, "zsh") {
		shellrc = filepath.Join(homeDir, ".zshrc")
	} else if strings.Contains(shell_config, "bash") {
		shellrc = filepath.Join(homeDir, ".bashrc")
	} else {
		theme.Warning(fmt.Sprintf("Could not detect shell config for %s. Please add ~/.pvm/current to your PATH manually.", shell_config))
		return
	}

	path_to_add := "export PATH=\"$HOME/.pvm/current:$PATH\""
	content, err := os.ReadFile(shellrc)
	if err != nil {
		// If the file doesn't exist, we can't check it.
		// We'll just print the message.
		theme.Warning(fmt.Sprintf("Could not read %s. Please add the following line to your shell configuration file:", shellrc))
		fmt.Println(path_to_add)
		theme.Success(fmt.Sprintf("Now using PHP %s", version))
		return
	}

	if !strings.Contains(string(content), "$HOME/.pvm/current") {
		theme.Warning("To finish setting up pvm, please add the following line to your shell configuration file:")
		fmt.Printf("  %s\n", path_to_add)

		theme.Info("You can do this automatically by running:")

		fmt.Printf("  echo '%s' | sudo tee -a %s > /dev/null && source %s\n", path_to_add, shellrc, shellrc)

		theme.Info("This will append the configuration with sudo and apply it immediately.")
		return
	}

	theme.Success(fmt.Sprintf("Now using PHP %s", version))
}

func useForWindows(args []string) {
	threadSafe := true
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

	versionFolderPath := filepath.Join(homeDir, ".pvm", "versions", selectedVersion.folder.Name())
	versionPath := filepath.Join(versionFolderPath, "php.exe")
	versionPathCGI := filepath.Join(versionFolderPath, "php-cgi.exe")
	composerPath := filepath.Join(versionFolderPath, "composer", "composer.phar")

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

	err = os.WriteFile(batPath, []byte(batCommandCGI), 0755)

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

	// create bat script for composer
	batCommandComposer := "@echo off \n"
	batCommandComposer = batCommandComposer + "set filepath=\"" + versionPath + "\"\n"
	batCommandComposer = batCommandComposer + "set composerpath=\"" + composerPath + "\"\n"
	batCommandComposer = batCommandComposer + "set arguments=%*\n"
	batCommandComposer = batCommandComposer + "%filepath% %composerpath% %arguments%\n"

	err = os.WriteFile(batPathComposer, []byte(batCommandComposer), 0755)

	if err != nil {
		log.Fatalln(err)
	}

	// create sh script for php
	shCommandComposer := "#!/bin/bash\n"
	shCommandComposer = shCommandComposer + "filepath=\"" + versionPath + "\"\n"
	shCommandComposer = shCommandComposer + "composerpath=\"" + composerPath + "\"\n"
	shCommandComposer = shCommandComposer + "\"$filepath\" \"$composerpath\" \"$@\""

	err = os.WriteFile(shPathComposer, []byte(shCommandComposer), 0755)

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
