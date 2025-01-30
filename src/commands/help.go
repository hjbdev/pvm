package commands

import (
	"fmt"
	"hjbdev/pvm/theme"
)

func Help(notFoundError bool) {
	theme.Title("pvm: PHP Version Manager")
	theme.Info("Version 1.3.1")

	if notFoundError {
		theme.Error("Error: Command not found.")
		fmt.Println()
	}

	fmt.Println("Usage:")
	fmt.Println("  pvm [command] [options]")
	fmt.Println()
	fmt.Println("Available Commands:")

	commands := map[string]string{
		"help":        "Display help information about pvm commands.",
		"install":     "Install a specific PHP version.",
		"uninstall":   "Uninstall a specific PHP version.",
		"list-remote": "List available PHP versions from remote repositories.",
		"list":        "List installed PHP versions.",
		"use":         "Switch to a specific installed PHP version.",
		"add":         "Add a custom PHP version source.",
		"remove":      "Remove a custom PHP version source.",
	}

	for cmd, desc := range commands {
		fmt.Printf("  %-12s %s\n", cmd, desc)
	}

	fmt.Println()
	fmt.Println("For detailed usage of a specific command, use:")
	fmt.Println("  pvm help [command]")
}
