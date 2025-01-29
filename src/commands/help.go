package commands

import (
	"fmt"
	"hjbdev/pvm/theme"
)

func Help(notFoundError bool) {
	theme.Title("pvm: PHP Version Manager")
	theme.Info("Version 1.3.0")

	if notFoundError {
		theme.Error("Command not found")
	}

	fmt.Println("Available Commands:")
	fmt.Println("    help")
	fmt.Println("    install")
	fmt.Println("    uninstall")
	fmt.Println("    list-remote")
	fmt.Println("    list")
	fmt.Println("    use")
}
