package nginxconf

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

var gConf Config

type Config struct {
	Conf   string `json:"conf"`
	Reload string `json:"reload"`
}

func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gConf)
	return err
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("nginx", Parse, Save)
}
