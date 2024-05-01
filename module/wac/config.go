package wac

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

var gConf Config

type Config struct {
	NginxPath string `json:"nginx_conf_path"`
	Reload    string `json:"reload"`
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
	config.RegisterConfig("wac", Parse, Save)
}
