package main

import (
	"flag"

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

	s := cmmns.New(configfile)
	// 在init之前，注册业务router

	// Init
	s.Init()
	// Run
	s.Run()

	defer s.Close()
	s.Wait()
}
