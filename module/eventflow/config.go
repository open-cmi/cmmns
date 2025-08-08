package eventflow

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

type Config struct {
	Relay    string            `json:"relay"`
	Method   string            `json:"method"`
	Channels map[string]string `json:"channels"`
}

var gConf Config

func GetConf() *Config {
	return &gConf
}

func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(gConf)
	return v
}

func init() {
	gConf.Channels = make(map[string]string)
	config.RegisterConfig("event_flow", Parse, Save)
}
