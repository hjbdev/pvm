package main

import (
	"hjbdev/pvm/commands"
	"os"
	"runtime"

	"hjbdev/pvm/theme"
)

func main() {
	args := os.Args[1:]

	runtimeOS := runtime.GOOS
	arch := runtime.GOARCH

	if runtimeOS != "windows" {
		theme.Error("pvm currently only works on Windows.")
		return
	}

	if arch != "amd64" {
		theme.Error("pvm currently only works on 64-bit systems.")
		return
	}

	if len(args) == 0 {
		commands.Help(false)
		return
	}

	var err error

	switch args[0] {
	case "help":
		commands.Help(false)
	case "update":
		err = commands.Update(args[1:])
	case "ls", "list":
		err = commands.List(args[1:])
	case "bin":
		err = commands.Bin()
	case "install", "i":
		err = commands.Install(args)
	case "use", "u":
		err = commands.Use(args[1:])
	case "extensions", "e":
		err = commands.Extensions(args[1:])
	default:
		commands.Help(true)
	}

	if err != nil {
		theme.Error(err.Error())
		os.Exit(1)
	}
}
