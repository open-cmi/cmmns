package cmmns

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"

	_ "github.com/open-cmi/cmmns/migration"
)

func Init(configFile string) error {
	err := config.Init(configFile)
	if err != nil {
		logger.Errorf("new config failed: %s\n", err.Error())
		return err
	}

	return nil
}

func Wait() {
	// 初始化后，等待信号
	sigs := make(chan os.Signal, 1)

	//注册要接收的信号，syscall.SIGINT:接收ctrl+c ,syscall.SIGTERM:程序退出
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigs
}

func Fini() {
	config.Save()
}
