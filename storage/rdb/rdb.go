package rdb

import (
	"fmt"

	"github.com/open-cmi/cmmns/config"

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

// Cache redis Cache
var Cache [MaxCache]*redis.Client

// Init db init
func Init() error {

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

// GetCache get cache
func GetCache(mod int) *redis.Client {
	return Cache[mod]
}
