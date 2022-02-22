package webserver

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/service/webserver/middleware"
)

type Service struct {
	Engine *gin.Engine
}

func New() *Service {
	return &Service{
		Engine: gin.New(),
	}
}

func (s *Service) Init() error {
	// init webserver
	err := middleware.Init(&moduleConfig.Middleware)
	if err != nil {
		logger.Error("middleware init failed")
		return err
	}

	middleware.DefaultMiddleware(s.Engine)
	middleware.SessionMiddleware(s.Engine)
	UnauthInit(s.Engine)
	middleware.AuthMiddleware(s.Engine)
	AuthInit(s.Engine)
	return nil
}

func (s *Service) Run() error {
	// unix sock api
	for _, srv := range moduleConfig.Server {
		if srv.Proto == "unix" {
			sockAddr := srv.Address
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
			go http.Serve(listener, s.Engine)
		} else if srv.Proto == "http" {
			go s.Engine.Run(srv.Address + ":" + strconv.Itoa(srv.Port))
		}
	}

	return nil
}

func (s *Service) Close() {

}
