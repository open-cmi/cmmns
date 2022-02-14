package middleware

import (
	"fmt"
	"strconv"

	"github.com/open-cmi/memstore"
	"github.com/topmyself/redistore"
)

type RedisStoreConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	DB     int    `json:"db"`
	Passwd string `json:"passwd"`
}

type MiddlewareConfig struct {
	SessionStore string           `json:"session_store"`
	RedisStore   RedisStoreConfig `json:"redis_store,omitempty"`
}

var middlewareConfig MiddlewareConfig

// Init init func
func Init(config *MiddlewareConfig) (err error) {
	if config.SessionStore == "memory" {
		memoryStore = memstore.NewMemStore([]byte("memorystore"),
			[]byte("enckey12341234567890123456789012"))
		storeType = "memory"
	} else {
		host := fmt.Sprintf("%s:%d", config.RedisStore.Host, config.RedisStore.Port)
		pass := config.RedisStore.Passwd
		redisStore, err = redistore.NewRediStoreWithDB(100, "tcp", host, pass, strconv.Itoa(config.RedisStore.DB))
		if err != nil {
			return err
		}

		redisStore.SetKeyPrefix("koa-sess-")
		redisStore.SetSerializer(redistore.JSONSerializer{})
	}

	return nil
}
