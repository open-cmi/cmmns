package main

import (
	"flag"

	"github.com/open-cmi/cmmns"
	"github.com/open-cmi/cmmns/scmd"
	"github.com/open-cmi/cmmns/service/ticker"
	"github.com/open-cmi/cmmns/service/webserver"
	"github.com/open-cmi/migrate"
)

var configfile string = ""

func main() {
	if migrate.TryRun("cmmns") {
		return
	}

	if scmd.TryRun() {
		return
	}

	flag.StringVar(&configfile, "config", configfile, "config file")
	flag.Parse()

	if configfile == "" {
		flag.Usage()
		return
	}

	err := cmmns.Init(configfile)
	if err != nil {
		return
	}
	defer cmmns.Fini()

	s := webserver.New()
	// Init
	s.Init()
	// Run
	s.Run()

	// run ticker service
	t := ticker.New()
	t.Init()
	t.Run()

	cmmns.Wait()
}
