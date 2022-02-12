package main

import (
	"flag"
	"net"
	"net/http"
	"os"

	"github.com/open-cmi/cmmns"
	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/middleware"
	"github.com/open-cmi/migrate"

	"github.com/gin-gonic/gin"
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

	err := cmmns.Init(configfile)
	if err != nil {
		logger.Logger.Error("service init failed\n")
		return
	}

	// init router
	r := gin.New()

	err = middleware.Init()
	if err != nil {
		logger.Logger.Error("middleware init failed")
		return
	}
	middleware.DefaultMiddleware(r)
	cmmns.NauthInit(r)
	middleware.AuthMiddleware(r)
	middleware.UserPermMiddleware(r)
	cmmns.AuthInit(r)

	// unix sock api
	const sockAddr = "/tmp/cmmns.sock"
	os.Remove(sockAddr)
	unixAddr, err := net.ResolveUnixAddr("unix", sockAddr)
	if err != nil {
		logger.Logger.Error(err.Error() + "\n")
		return
	}

	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		logger.Logger.Error("listening error: %s\n", err.Error())
	}
	logger.Logger.Debug("listening unix socket: %s\n", sockAddr)
	go http.Serve(listener, r)

	r.Run(":30000")
}
