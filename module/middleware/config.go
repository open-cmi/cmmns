package middleware

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/rdb"
	"github.com/open-cmi/cmmns/service/business"
	"github.com/open-cmi/memstore"
	"github.com/topmyself/redistore"
)

type Config struct {
	Store  string `json:"store"`
	MaxAge int    `json:"max_age"`
}

var gConf Config

func Init() error {
	var err error
	rdbConf := rdb.GetConf()
	if gConf.Store == "memory" {
		memoryStore = memstore.NewMemStore([]byte("memorystore"),
			[]byte("enckey12341234567890123456789012"))
		memoryStore.MaxAge(gConf.MaxAge)
	} else if gConf.Store == "redis" {
		host := fmt.Sprintf("%s:%d", rdbConf.Host, rdbConf.Port)
		redisStore, err = redistore.NewRediStoreWithDB(100, "tcp", host, rdbConf.Password, "2")
		if err != nil {
			logger.Error("redis store new failed\n")
			return err
		}

		redisStore.SetKeyPrefix("koa-sess-")
		redisStore.SetSerializer(redistore.JSONSerializer{})
		redisStore.SetMaxAge(gConf.MaxAge)
	} else {
		return errors.New("middleware store type not supported")
	}

	return nil
}

func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}
	if gConf.MaxAge == 0 {
		gConf.MaxAge = 3600
	}
	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	// default config
	gConf.Store = "memory"

	config.RegisterConfig("middleware", Parse, Save)
	business.Register("middleware", business.DefaultPriority, Init)
}
