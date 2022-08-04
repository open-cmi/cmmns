package user

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

// Config smtp
type Config struct {
	From       string `json:"from"`
	SMTPServer string `json:"smtp_server"`
	User       string `json:"user"`
	Password   string `json:"password"`
	SMTPHost   string `json:"smtp_host"`
	Domain     string `json:"domain"`
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
	config.RegisterConfig("smtp", Init, Save)
}
