package user2

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

// Config smtp
type Config struct {
	ActivateURL string `json:"activate_url"`
}

func Init(confmsg json.RawMessage) error {
	err := json.Unmarshal(confmsg, &gConf)
	return err
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

var gConf Config

func init() {
	config.RegisterConfig("user2", Init, Save)
}
