package main

import (
	"flag"

	"github.com/open-cmi/cmmns"
)

var configfile string = ""

func main() {
	if cmmns.TryRunScmd() {
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

	cmmns.Run()

	cmmns.Wait()
}
