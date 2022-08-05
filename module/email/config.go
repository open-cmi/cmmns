package email

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type Config struct {
	Server   string
	Port     int
	User     string
	Password string
}

var gConf Config

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	return err
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("smtp", Init, Save)
}
