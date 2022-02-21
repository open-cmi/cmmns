package webserver

import (
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
	err := Init(s.Engine)
	if err != nil {
		logger.Error("middleware init failed")
		return err
	}

	middleware.SessionMiddleware(s.Engine)
	UnauthInit(s.Engine)
	middleware.AuthMiddleware(s.Engine)
	AuthInit(s.Engine)

	return nil
}

func (s *Service) Run() error {
	// unix sock api
	Run(s.Engine)
	return nil
}

func (s *Service) Close() {
	Close()
}
