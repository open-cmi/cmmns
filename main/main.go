package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-cmi/cmmns"
	"github.com/open-cmi/migrate"
)

var configfile string = ""

func main() {
	if migrate.TryRun("cmmns") {
		return
	}

	flag.StringVar(&configfile, "config", configfile, "config file")
	flag.Parse()

	if configfile == "" {
		flag.Usage()
		return
	}
	// 读取配置文件
	/*
		conf, err := config.Init(configfile)
		if err != nil {
			fmt.Printf("config init failed: %s\n", err.Error())
			return
		}*/

	s := cmmns.New(configfile)
	// 在init之前，注册业务router

	// Init
	s.Init()
	// Run
	s.Run()

	defer s.Close()
	// 初始化后，等待信号
	sigs := make(chan os.Signal, 1)

	//注册要接收的信号，syscall.SIGINT:接收ctrl+c ,syscall.SIGTERM:程序退出
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigs
}
