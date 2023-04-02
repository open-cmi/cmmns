package agent

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type Config struct {
	CheckStatus bool `json:"enable,omitempty"`
}

var gConf Config

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	return err
}

func Save() json.RawMessage {
	b, _ := json.Marshal(gConf)
	return b
}

func init() {
	gConf.CheckStatus = false
	config.RegisterConfig("agent", Init, Save)
}
