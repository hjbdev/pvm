package main

import (
	"hjbdev/pvm/commands"
	"os"
	"runtime"

	"github.com/fatih/color"
)

func main() {
	args := os.Args[1:]

	os := runtime.GOOS

	if os != "windows" {
		color.Red("pvm currently only works on Windows.")
		color.Red("This OS is not supported and may not function correctly.")
		color.White("")
	}

	if len(args) == 0 {
		commands.Help(false)
		return
	}

	switch args[0] {
	case "help":
		commands.Help(false)
	// case "list":
	// 	commands.List()
	case "install":
		commands.Install(args)
	case "discover":
		commands.Discover(args[1:])
	default:
		commands.Help(true)
	}
}
