package commands

import (
	"fmt"

	"github.com/fatih/color"
)

func Help(notFoundError bool) {
	fmt.Println()
	color.Blue("pvm: PHP Version Manager")
	fmt.Println()

	if notFoundError {
		color.Red("Command not found")
		fmt.Println()
	}

	fmt.Println("Available Commands:")
	fmt.Println("    help")
	fmt.Println("    list")
	fmt.Println("    install")
	fmt.Println("    discover")
	fmt.Println()
}
