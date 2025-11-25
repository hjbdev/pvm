package main

import (
	"hjbdev/pvm/commands"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		commands.Help(false)
		return
	}

	switch args[0] {
	case "help":
		commands.Help(false)
	case "ls", "list":
		commands.List()
	case "ls remote", "ls-remote", "list remote", "list-remote":
		commands.ListRemote()
	case "path":
		commands.Path()
	case "install":
		commands.Install(args)
	case "use":
		commands.Use(args[1:])
	case "uninstall":
		commands.Uninstall(args)
	default:
		commands.Help(true)
	}
}
