package cmmns

import (
	"github.com/open-cmi/cmmns/essential/subcommands"

	_ "github.com/open-cmi/cmmns/api"
	_ "github.com/open-cmi/cmmns/commands"
	_ "github.com/open-cmi/cmmns/migration"
	_ "github.com/open-cmi/cmmns/module"
)

func Main() error {
	return subcommands.Run()
}
