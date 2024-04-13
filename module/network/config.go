package network

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

var gConf Config

type DevConfig struct {
	DHCP         bool   `json:"dhcp"`
	Address      string `json:"address,omitempty"`
	Netmask      string `json:"netmask,omitempty"`
	Gateway      string `json:"gateway,omitempty"`
	PreferredDNS string `json:"preferred_dns,omitempty"`
	AlternateDNS string `json:"alternate_dns,omitempty"`
}

type Config struct {
	Engine   string               `json:"engine,omitempty"`
	ConfFile string               `json:"conf_file"`
	Devices  map[string]DevConfig `json:"devices"`
}

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	if gConf.Engine == "" {
		gConf.Engine = "netplan"
	}

	err = NetworkApply(&gConf)

	return err
}

func (c *Config) Save() {
	config.Save()
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("network", Init, Save)
}
