package run

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/subcommands"
	"github.com/open-cmi/cmmns/service/business"
	"github.com/open-cmi/cmmns/service/ticker"
	"github.com/open-cmi/cmmns/service/webserver"
)

type RunCommand struct {
	ConfigFile string
}

func Wait() {
	// 初始化后，等待信号
	sigs := make(chan os.Signal, 1)

	//注册要接收的信号，syscall.SIGINT:接收ctrl+c ,syscall.SIGTERM:程序退出
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigs
}

func Fini() {
}

func (c *RunCommand) Synopsis() string {
	return "run main program"
}

func (c *RunCommand) Run() error {

	if os.Args[1] == "run" {
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		runCmd.StringVar(&c.ConfigFile, "config", c.ConfigFile, "config file")

		err := runCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
		if c.ConfigFile == "" {
			runCmd.Usage()
			return errors.New("config file must not be empty")
		}
	} else {
		flag.StringVar(&c.ConfigFile, "config", c.ConfigFile, "config file")
		flag.Parse()

		if c.ConfigFile == "" {
			flag.Usage()
			return errors.New("config file must not be empty")
		}
	}

	err := config.Init(c.ConfigFile)
	if err != nil {
		logger.Errorf("config init failed: %s\n", err.Error())
		return err
	}
	err = business.Init()
	if err != nil {
		logger.Errorf("%s\n", err.Error())
		return err
	}

	defer Fini()

	// start web service
	s := webserver.New()

	// Init
	err = s.Init()
	if err != nil {
		return err
	}

	// Run
	err = s.Run()
	if err != nil {
		return err
	}

	// run ticker service
	t := ticker.New()
	err = t.Init()
	if err != nil {
		return err
	}
	err = t.Run()
	if err != nil {
		return nil
	}
	Wait()
	return nil
}

func init() {
	subcommands.RegisterCommand("run", &RunCommand{})
	subcommands.DefaultCommand(&RunCommand{})
}
