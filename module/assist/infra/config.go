package infra

import "github.com/open-cmi/cmmns/essential/config"

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

var moduleConfig Config

func init() {
	config.RegisterConfig("assist", &moduleConfig)
}
