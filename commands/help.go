package commands

import (
	"fmt"
	"hjbdev/pvm/theme"
)

func Help(notFoundError bool) {
	theme.Title("pvm: PHP Version Manager")
	theme.Info("Version 1.0")

	if notFoundError {
		theme.Error("Command not found")
	}

	fmt.Println("Available Commands:")
	fmt.Println("    help")
	fmt.Println("    install")
	fmt.Println("    list")
	fmt.Println("    path")
	fmt.Println("    use")
}
