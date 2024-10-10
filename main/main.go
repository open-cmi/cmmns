package main

import (
	"fmt"
	"os"

	"github.com/open-cmi/cmmns"

	_ "github.com/open-cmi/cmmns/internal/translation"
)

func main() {
	err := cmmns.Main()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
