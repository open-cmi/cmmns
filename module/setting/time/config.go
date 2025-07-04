package time

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

type Config struct {
	Manage     string   `json:"manage"`
	NTPServers []string `json:"ntp_servers"`
}

var gConf Config

func ParseConfig(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	return err
}

func SaveConfig() json.RawMessage {
	v, _ := json.Marshal(gConf)
	return v
}

func init() {
	config.RegisterConfig("time", ParseConfig, SaveConfig)
}
