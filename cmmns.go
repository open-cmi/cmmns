package cmmns

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/service/ticker"
	"github.com/open-cmi/cmmns/service/webserver"

	_ "github.com/open-cmi/cmmns/api"
	_ "github.com/open-cmi/cmmns/internal/translation"
	_ "github.com/open-cmi/cmmns/migration"
	_ "github.com/open-cmi/cmmns/module"
)

type Option struct {
	WebServiceEnable  bool
	TickServiceEnable bool
}

func Init(configFile string) error {
	err := config.Init(configFile)
	if err != nil {
		logger.Errorf("config init failed: %s\n", err.Error())
		return err
	}

	return nil
}

func Run(opt *Option) error {
	var count int
	if opt.WebServiceEnable {
		s := webserver.New()
		// Init
		s.Init()
		// Run
		s.Run()
		count++
	}

	if opt.TickServiceEnable {
		// run ticker service
		t := ticker.New()
		t.Init()
		t.Run()
		count++
	}

	if count > 1 {
		return nil
	}

	return errors.New("at least one service enabled")
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
