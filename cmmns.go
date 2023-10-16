package cmmns

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/open-cmi/cmmns/essential/subcommands"
	"github.com/open-cmi/migrate"

	_ "github.com/open-cmi/cmmns/api"
	_ "github.com/open-cmi/cmmns/commands/run"
	_ "github.com/open-cmi/cmmns/internal/translation"
	_ "github.com/open-cmi/cmmns/migration"
	_ "github.com/open-cmi/cmmns/module"
)

func Run() error {

	if len(os.Args) < 2 {
		flag.Usage()
		return errors.New("program args is incorrect")
	}

	subcmd := os.Args[1]
	if !strings.HasPrefix(subcmd, "-") {
		if migrate.TryRun() {
			return nil
		}
	}

	return subcommands.Run()
}
