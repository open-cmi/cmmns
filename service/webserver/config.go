package webserver

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type Server struct {
	Address  string `json:"address"`
	Port     int    `json:"port,omitempty"`
	Proto    string `json:"proto"`
	CertFile string `json:"cert,omitempty"`
	KeyFile  string `json:"key,omitempty"`
}

type Config struct {
	Server []Server `json:"server"`
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
	// default config

	gConf.Server = append(gConf.Server, Server{
		Address: "127.0.0.1",
		Port:    30000,
		Proto:   "http",
	}, Server{
		Address: "/tmp/cmmns.sock",
		Proto:   "unix",
	})

	config.RegisterConfig("webserver", Init, Save)
}
