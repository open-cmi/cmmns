package rdb

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/cmmns/essential/config"
)

// Cache redis Cache
var Cache map[string]*redis.Client

var modules map[string]int = make(map[string]int)

// Config cache config
type Config struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
}

var moduleConfig Config

// Init db init
func (c *Config) Init() error {
	cachehost := c.Host
	cacheport := c.Port
	cachepassword := c.Password

	Cache = make(map[string]*redis.Client)

	for module, db := range modules {
		Cache[module] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cachehost, cacheport),
			Password: cachepassword,
			DB:       db,
		})
	}

	return nil
}

// GetCache get cache
func GetCache(module string) *redis.Client {
	return Cache[module]
}

func Register(module string, db int) error {
	_, found := modules[module]
	if found {
		return errors.New("module has been registered")
	}
	modules[module] = db
	return nil
}

func init() {
	config.RegisterConfig("redis", &moduleConfig)
}
