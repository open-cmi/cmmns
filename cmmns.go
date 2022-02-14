package cmmns

import (
	"github.com/gin-gonic/gin"
	_ "github.com/open-cmi/cmmns/component"
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/api/middleware"
	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/storage"
	"github.com/open-cmi/cmmns/essential/ticker"
)

type APIGroup func(e *gin.Engine)

type Service struct {
	Engine *gin.Engine
}

func New(configFile string) *Service {
	err := config.Init(configFile)
	if err != nil {
		logger.Errorf("new config failed: %s\n", err.Error())
		return nil
	}

	return &Service{
		Engine: gin.New(),
	}
}

func (s *Service) Init() error {
	logger.Init()

	// 在这里会调用各个模块的配置函数

	// init router
	err := api.Init(s.Engine)
	if err != nil {
		logger.Error("middleware init failed")
		return err
	}

	api.UnauthInit(s.Engine)
	middleware.AuthMiddleware(s.Engine)
	middleware.UserPermMiddleware(s.Engine)
	api.AuthInit(s.Engine)
	ticker.Init()
	storage.Init()
	return nil
}

func (s *Service) Run() error {
	// unix sock api
	api.Run(s.Engine)
	ticker.Run()
	return nil
}

func (s *Service) Close() {
	ticker.Close()
	api.Close()
	config.Save()
}
