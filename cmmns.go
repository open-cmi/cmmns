package cmmns

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/cmmns/log"
	"github.com/open-cmi/cmmns/router"
	"github.com/open-cmi/cmmns/storage"
	"github.com/open-cmi/cmmns/ticker"
)

// Init service Init
func Init(configfile string) error {
	log.Init()

	// 配置文件的配置先确定在这里
	err := config.Init(configfile)
	if err != nil {
		log.Logger.Printf(log.Error, "%s\n", err.Error())
		return err
	}

	storage.Init()
	ticker.Init()
	return nil
}

// AuthInit auth init
func AuthInit(e *gin.Engine) {
	router.AuthInit(e)
}

// NauthInit no auth init
func NauthInit(e *gin.Engine) {
	router.NauthInit(e)
}
