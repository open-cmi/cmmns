package middleware

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/memstore"
	"github.com/topmyself/redistore"
)

type RedisStoreConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	DB     int    `json:"db"`
	Passwd string `json:"passwd"`
}

type Config struct {
	StoreType  string           `json:"store_type"`
	RedisStore RedisStoreConfig `json:"redis_store,omitempty"`
}

var gConf Config

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	if gConf.StoreType == "memory" {
		memoryStore = memstore.NewMemStore([]byte("memorystore"),
			[]byte("enckey12341234567890123456789012"))
		storeType = "memory"
	} else {
		host := fmt.Sprintf("%s:%d", gConf.RedisStore.Host, gConf.RedisStore.Port)
		pass := gConf.RedisStore.Passwd
		redisStore, err = redistore.NewRediStoreWithDB(100, "tcp", host, pass, strconv.Itoa(gConf.RedisStore.DB))
		if err != nil {
			return err
		}

		redisStore.SetKeyPrefix("koa-sess-")
		redisStore.SetSerializer(redistore.JSONSerializer{})
	}

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	// default config
	gConf.StoreType = "memory"

	config.RegisterConfig("session_middleware", Init, Save)
}
