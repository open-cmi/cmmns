package db

import (
	"database/sql"
	"fmt"

	"github.com/open-cmi/cmmns/config"

	"github.com/open-cmi/goutils/database"
	"github.com/open-cmi/goutils/database/dbsql"

	"github.com/go-redis/redis/v8"
)

const (
	// UserCache user cache
	UserCache = 0
	// TaskCache task cache
	TaskCache = 1

	// AgentCache agent cache
	AgentCache = 2

	// MaxCache max cache
	MaxCache = 3
)

// DB sql db
var DB *sql.DB

// Cache redis Cache
var Cache [MaxCache]*redis.Client

// Init db init
func Init() error {
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
		return err
	}
	DB = dbi

	rdb := config.GetConfig().Rdb
	cachehost := rdb.Host
	cacheport := rdb.Port
	cachepassword := rdb.Password

	Cache[UserCache] = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cachehost, cacheport),
		Password: cachepassword,
		DB:       1,
	})
	Cache[TaskCache] = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cachehost, cacheport),
		Password: cachepassword,
		DB:       1,
	})
	Cache[AgentCache] = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cachehost, cacheport),
		Password: cachepassword,
		DB:       1,
	})
	return nil
}

// GetDB get db
func GetDB() *sql.DB {
	return DB
}

// GetCache get cache
func GetCache(mod int) *redis.Client {
	return Cache[mod]
}
