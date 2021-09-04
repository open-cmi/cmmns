package cmmns

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/cmmns/db"
	"github.com/open-cmi/cmmns/router"
	"github.com/open-cmi/cmmns/startup"

	"github.com/open-cmi/goutils/database"
	"github.com/open-cmi/goutils/database/dbsql"

	"github.com/go-redis/redis/v8"
)

// Init service Init
func Init(configfile string) error {
	// 配置文件的配置先确定在这里
	err := config.Init(configfile)
	if err != nil {
		fmt.Printf("read config file error %s\n", err.Error())
		return err
	}
	var dbconf database.Config
	model := config.GetConfig().Model
	dbconf.Type = model.Type
	dbconf.Host = model.Host
	dbconf.Port = model.Port
	dbconf.User = model.User
	dbconf.Password = model.Password
	dbconf.Database = model.Database

	dbi, err := dbsql.SQLInit(&dbconf)
	if err != nil {
		fmt.Printf("db sql init failed\n")
		return err
	}
	db.DB = dbi

	rdb := config.GetConfig().Rdb
	cachehost := rdb.Host
	cacheport := rdb.Port
	cachepassword := rdb.Password

	cache := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cachehost, cacheport),
		Password: cachepassword,
		DB:       1,
	})
	db.Cache = cache

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
