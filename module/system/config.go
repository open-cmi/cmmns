package system

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

type Config struct {
	MonitorUsage bool `json:"monitor_usage"`
}

var gConf Config

func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(&gConf)
	return v
}

func init() {
	config.RegisterConfig("system", Parse, Save)
}
