package main

import (
	"fmt"
	"os"

	_ "github.com/open-cmi/cmmns/internal/commands"
	_ "github.com/open-cmi/cmmns/internal/translation"

	_ "github.com/open-cmi/cmmns"
	"github.com/open-cmi/gobase"
)

func main() {
	err := gobase.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
