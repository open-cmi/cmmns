package main

import (
	"fmt"
	"os"

	_ "github.com/open-cmi/cmmns/internal/commands"
	_ "github.com/open-cmi/cmmns/internal/translation"

	"github.com/open-cmi/cmmns"
)

func main() {
	err := cmmns.Main()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
