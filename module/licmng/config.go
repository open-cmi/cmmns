package licmng

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

var gConf Config

type Config struct {
	Enable      bool   `json:"enable"`
	PrivateFile string `json:"private_file"`
}

func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gConf)
	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(gConf)
	return v
}

func init() {
	config.RegisterConfig("licmng", Parse, Save)
}
