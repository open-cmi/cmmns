package subcommands

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type Command interface {
	Run() error
	Synopsis() string
}

var cmdmapping map[string]Command = make(map[string]Command)
var RootCommand Command

func DefaultCommand(cmd Command) {
	RootCommand = cmd
}

func RegisterCommand(key string, cmd Command) error {
	_, found := cmdmapping[key]
	if found {
		return errors.New("command has been registered")
	}

	cmdmapping[key] = cmd

	return nil
}

func Usage() {
	prog := path.Base(os.Args[0])
	fmt.Printf("Usage: %s <subcommand> <subcommand args>\n\n", prog)

	fmt.Printf("Subcommands:\n")
	for name, command := range cmdmapping {
		fmt.Printf("\t%-15s\t%s\n", name, command.Synopsis())
	}
	return
}

func Run() error {
	if len(os.Args) < 2 {
		Usage()
		return errors.New("subcommand cli args is too short")
	}

	var err error
	key := os.Args[1]
	cmd, ok := cmdmapping[key]
	if ok {
		err = cmd.Run()
	} else {
		err = RootCommand.Run()
	}

	return err
}
