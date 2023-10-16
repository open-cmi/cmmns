package subcommands

import (
	"errors"
	"os"
)

type Command interface {
	Run() error
}

var cmdmapping map[string]Command = make(map[string]Command)
var RootCommand Command

func RegisterCommand(key string, isRootCmd bool, cmd Command) error {
	_, found := cmdmapping[key]
	if found {
		return errors.New("command has been registered")
	}

	cmdmapping[key] = cmd
	if isRootCmd {
		RootCommand = cmd
	}
	return nil
}

func Run() error {
	if len(os.Args) < 2 {
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
