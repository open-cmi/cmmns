package main

import (
	"flag"

	"github.com/open-cmi/cmmns"
)

var configfile string = ""

func main() {
	if cmmns.TryRunScmd("cmmns") {
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

	var opt cmmns.Option
	opt.WebServiceEnable = true
	opt.TickServiceEnable = true
	cmmns.Run(&opt)

	cmmns.Wait()
}
