package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/open-cmi/cmmns"
	"github.com/open-cmi/cmmns/middleware"
	"github.com/open-cmi/cmmns/migration"
	"github.com/open-cmi/migrate"

	"github.com/gin-gonic/gin"
)

var configfile string = ""

func main() {
	if len(os.Args) > 1 && migrate.IsMigrateCommand(os.Args[1]) {
		migration.Migrate()
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
		fmt.Printf("service init failed\n")
		return
	}
	// init router
	r := gin.Default()

	err = middleware.Init()
	if err != nil {
		fmt.Printf("middleware init failed")
		return
	}

	middleware.AuthMiddleware(r)
	cmmns.NauthInit(r)
	middleware.UserPermMiddleware(r)
	cmmns.AuthInit(r)

	// unix sock api
	const sockAddr = "/tmp/cmmns.sock"
	os.Remove(sockAddr)
	unixAddr, err := net.ResolveUnixAddr("unix", sockAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		fmt.Println("listening error:", err)
	}
	fmt.Printf("listening unix socket: %s\n", sockAddr)
	go http.Serve(listener, r)

	r.Run(":30000")
}
