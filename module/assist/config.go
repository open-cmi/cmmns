package assist

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type RemoteService struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	LocalIP    string `json:"local_ip"`
	LocalPort  uint16 `json:"local_port"`
	RemotePort uint16 `json:"remote_port"`
}

type Config struct {
	ServerAddr string          `json:"server_addr"`
	ServerPort uint16          `json:"server_port"`
	Token      string          `json:"token,omitempty"`
	Service    []RemoteService `json:"services,omitempty"`
}

var gConf Config

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	return err
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("assist", Init, Save)
}
