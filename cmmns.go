package cmmns

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/cmmns/db"
	"github.com/open-cmi/cmmns/router"
	"github.com/open-cmi/cmmns/startup"
)

// Init service Init
func Init(configfile string) error {
	// 配置文件的配置先确定在这里
	err := config.Init(configfile)
	if err != nil {
		fmt.Printf("read config file error %s\n", err.Error())
		return err
	}

	db.Init()

	startup.Init()
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
