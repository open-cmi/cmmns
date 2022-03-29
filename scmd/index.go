package scmd

import (
	"errors"
	"os"
)

type Command interface {
	Run()
}

var cmdmapping map[string]Command = make(map[string]Command)

func RegisterCommand(key string, cmd Command) error {
	_, found := cmdmapping[key]
	if found {
		return errors.New("command has been registered")
	}

	cmdmapping[key] = cmd
	return nil
}

func TryRun() bool {
	if len(os.Args) <= 1 {
		return false
	}
	key := os.Args[1]
	cmd, ok := cmdmapping[key]
	if ok {
		cmd.Run()
		return true
	}
	return false
}

func Run() {
	if len(os.Args) <= 1 {
		return
	}
	key := os.Args[1]
	cmd, ok := cmdmapping[key]
	if ok {
		cmd.Run()
	}
}
