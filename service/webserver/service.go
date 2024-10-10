package webserver

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/middleware"
	"github.com/open-cmi/cmmns/pkg/eyas"
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

	middleware.DefaultMiddleware(s.Engine)
	//middleware.WACMiddleware(s.Engine)

	workDir := eyas.GetRootPath()
	dir := fmt.Sprintf("%s/static/", workDir)
	s.Engine.Static("/api-static/", dir)
	middleware.SessionMiddleware(s.Engine)
	middleware.JWTMiddleware(s.Engine)
	UnauthInit(s.Engine)
	middleware.AuthMiddleware(s.Engine)
	AuthInit(s.Engine)
	return nil
}

func (s *Service) Run() error {
	// unix sock api
	for _, srv := range gConf.Server {
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
		} else if srv.Proto == "https" {
			rp := eyas.GetRootPath()
			var certFile string
			if strings.HasPrefix(srv.CertFile, ".") {
				certFile = filepath.Join(rp, srv.CertFile)
			} else {
				certFile = srv.CertFile
			}
			var keyFile string
			if strings.HasPrefix(srv.KeyFile, ".") {
				keyFile = filepath.Join(rp, srv.KeyFile)
			} else {
				keyFile = srv.KeyFile
			}
			logger.Debugf("run tls %s:%d, cert %s, key %s",
				srv.Address, srv.Port, certFile, keyFile)
			go s.Engine.RunTLS(srv.Address+":"+strconv.Itoa(srv.Port),
				certFile, keyFile)
		}
	}

	return nil
}

func (s *Service) Close() {

}
