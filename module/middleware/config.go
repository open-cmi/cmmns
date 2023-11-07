package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/rdb"
	"github.com/open-cmi/memstore"
	"github.com/topmyself/redistore"
)

type Config struct {
	StoreType string `json:"store_type"`
}

var gConf Config

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}
	rdbConf := rdb.GetConf()
	if gConf.StoreType == "memory" {
		memoryStore = memstore.NewMemStore([]byte("memorystore"),
			[]byte("enckey12341234567890123456789012"))
		storeType = "memory"
	} else {
		host := fmt.Sprintf("%s:%d", rdbConf.Host, rdbConf.Port)
		redisStore, err = redistore.NewRediStoreWithDB(100, "tcp", host, rdbConf.Password, "2")
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

	config.RegisterConfig("middleware", Init, Save)
}
