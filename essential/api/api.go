package api

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/api/middleware"
	"github.com/open-cmi/cmmns/essential/logger"
)

func Init(e *gin.Engine) error {
	err := middleware.Init(&moduleConfig.Middleware)
	middleware.DefaultMiddleware(e)
	return err
}

func Run(r *gin.Engine) error {
	sockAddr := moduleConfig.UnixPath //"/tmp/cmmns.sock"
	os.Remove(sockAddr)
	unixAddr, err := net.ResolveUnixAddr("unix", sockAddr)
	if err != nil {
		logger.Error(err.Error() + "\n")
		return err
	}

	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		logger.Errorf("listening error: %s\n", err.Error())
	}

	logger.Debugf("listening unix socket: %s\n", sockAddr)
	go http.Serve(listener, r)

	go r.Run(moduleConfig.Listen + ":" + strconv.Itoa(moduleConfig.Port))
	return nil
}

func Close() {

}
