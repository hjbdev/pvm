package commands

import (
	"fmt"
	"hjbdev/pvm/theme"
	"os"

	"github.com/fatih/color"
)

var version = "dev"

func Help(notFoundError bool) {
	theme.Title("pvm: PHP Version Manager")
	theme.Info(fmt.Sprintf("Version %s", version))

	// Check for updates in background (best-effort). If env PVM_AUTO_UPDATE=1, run installer.
	if latest, newer, err := CheckForUpdate(version); err == nil {
		// show latest always
		theme.Info(fmt.Sprintf("Latest %s", latest))
		if newer {
			theme.Info("A newer version is available.")
			if os.Getenv("PVM_AUTO_UPDATE") == "1" {
				theme.Info("Auto-update enabled. Running installer...")
				if err := InstallLatest(true); err != nil {
					theme.Error(fmt.Sprintf("Auto-install failed: %v", err))
				}
			} else {
				theme.Info(fmt.Sprintf("Run the installer: irm https://pvm.hjb.dev/install.ps1 | iex"))
			}
		}
	} else {
		// non-fatal: hide network error
	}

	if notFoundError {
		theme.Error("Command not found")
	}

	fmt.Println("Available Commands:")
	printHelpCommand("extensions <list|ls|enable|disable> [extension[,extension...]]", "e")
	printHelpCommand("help")
	printHelpCommand("install", "i")
	printHelpCommand("update", "")
	printHelpCommand("list [remote]", "ls")
	printHelpCommand("bin")
	printHelpCommand("use <version>", "u")
}

func printHelpCommand(command string, aliases ...string) {
	line := "    " + command
	if len(aliases) == 0 {
		fmt.Println(line)
		return
	}

	fmt.Println(line + " " + color.HiBlackString("(alias: %s)", aliases[0]))
}
