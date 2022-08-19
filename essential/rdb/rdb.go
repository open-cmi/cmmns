package rdb

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/essential/config"
)

// clients redis clients
var clients map[string]*redis.Client

var modules map[string]int = make(map[string]int)

// Config cache config
type Config struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
}

var gConf Config

// GetClient get client
func GetClient(module string) *redis.Client {
	return clients[module]
}

func Register(module string, db int) error {
	_, found := modules[module]
	if found {
		return errors.New("module has been registered")
	}
	modules[module] = db
	return nil
}

// Init db init
func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	cachehost := gConf.Host
	cacheport := gConf.Port
	cachepassword := gConf.Password

	clients = make(map[string]*redis.Client)

	for module, db := range modules {
		clients[module] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cachehost, cacheport),
			Password: cachepassword,
			DB:       db,
		})
	}

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	gConf.Host = "127.0.0.1"
	gConf.Port = 25431
	config.RegisterConfig("redis", Init, Save)

	Register("public", def.RDBPublic)
}
