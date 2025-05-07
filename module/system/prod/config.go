package prod

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type Config struct {
	Name   string `json:"name"`
	Footer string `json:"footer"`
}

var gConf Config

func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gConf)
	return err
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("prod", Parse, Save)
}
