package current

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/subcommands"
	"github.com/open-cmi/cmmns/module/license"
)

type activateCommand struct {
}

var configfile string
var serialString string
var prodString string

func (c *activateCommand) Synopsis() string {
	return "activate product"
}

func (c *activateCommand) Run() error {
	currentCmd := flag.NewFlagSet("activate", flag.ExitOnError)
	currentCmd.StringVar(&serialString, "serial", serialString, "product serial")
	currentCmd.StringVar(&prodString, "product", prodString, "product name")
	currentCmd.StringVar(&configfile, "config", configfile, "config file")

	err := currentCmd.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	if configfile == "" {
		return errors.New("config file must not be empty")
	}

	err = config.Init(configfile)
	if err != nil {
		fmt.Printf("init config failed: %s\n", err.Error())
		return err
	}

	if serialString == "" || prodString == "" {
		currentCmd.Usage()
		return errors.New("product or serial is empty")
	}

	err = license.SetProductSerial(serialString, prodString)
	if err != nil {
		fmt.Printf("activate failed:%s\n", err.Error())
		return err
	}
	fmt.Println("verify serial successfully")
	fmt.Println("please restart service or wait for several minutes")
	return nil
}

func init() {
	subcommands.RegisterCommand("activate", &activateCommand{})
}
