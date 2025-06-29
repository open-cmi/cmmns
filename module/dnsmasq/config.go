package dnsmasq

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

type Config struct {
	Enable     bool   `json:"enable"`
	ConfigFile string `json:"config_file"`
	Start      string `json:"start"`
	Stop       string `json:"stop"`
	Restart    string `json:"restart"`
}

var gConf Config

func Parse(confmsg json.RawMessage) error {
	err := json.Unmarshal(confmsg, &gConf)

	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(gConf)
	return v
}

func init() {
	gConf.ConfigFile = "/etc/dnsmasq.conf"
	config.RegisterConfig("dnsmasq", Parse, Save)
}
