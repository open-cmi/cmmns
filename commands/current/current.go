package current

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/essential/subcommands"
	"github.com/open-cmi/migrate"
)

type CurrentCommand struct {
}

var configfile string

func (c *CurrentCommand) Synopsis() string {
	return "print current migrations in database"
}

func (c *CurrentCommand) Run() error {
	currentCmd := flag.NewFlagSet("current", flag.ExitOnError)
	currentCmd.StringVar(&configfile, "config", configfile, "config file")

	err := currentCmd.Parse(os.Args[2:])
	if err != nil {
		return err
	}
	if configfile == "" {
		currentCmd.Usage()
		return errors.New("config file must not be empty")
	}

	err = config.Init(configfile)
	if err != nil {
		fmt.Printf("init config failed: %s\n", err.Error())
		return err
	}

	opt := migrate.NewCurrentOpt(sqldb.GetConfDB())

	err = opt.Run()
	return err
}

func init() {
	subcommands.RegisterCommand("current", &CurrentCommand{})
}
