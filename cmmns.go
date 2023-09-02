package cmmns

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/scmd"
	"github.com/open-cmi/cmmns/service/business"
	"github.com/open-cmi/cmmns/service/ticker"
	"github.com/open-cmi/cmmns/service/webserver"
	"github.com/open-cmi/migrate"

	_ "github.com/open-cmi/cmmns/api"
	_ "github.com/open-cmi/cmmns/internal/translation"
	_ "github.com/open-cmi/cmmns/migration"
	_ "github.com/open-cmi/cmmns/module"
)

func TryRunScmd() bool {
	if migrate.TryRun() {
		return true
	}

	if scmd.TryRun() {
		return true
	}
	return false
}

func Init(configFile string) error {
	err := config.Init(configFile)
	if err != nil {
		logger.Errorf("config init failed: %s\n", err.Error())
		return err
	}
	err = business.Init()
	if err != nil {
		logger.Errorf("%s\n", err.Error())
		return err
	}

	return nil
}

func Run() error {
	// start web service
	s := webserver.New()
	// Init
	err := s.Init()
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

	return err
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
