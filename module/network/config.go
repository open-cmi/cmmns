package network

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type Config struct {
	Dev          string `json:"dev"`
	DHCP         bool   `json:"dhcp"`
	ConfFile     string `json:"conf_file"`
	Address      string `json:"address,omitempty"`
	Netmask      string `json:"netmask,omitempty"`
	Gateway      string `json:"gateway,omitempty"`
	MainDNS      string `json:"main_dns,omitempty"`
	SecondaryDNS string `json:"secondary_dns,omitempty"`
}

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	if gConf.Dev == "" {
		return nil
	}

	var msg ConfigMsg
	if gConf.DHCP {
		msg.Mode = "dhcp"
	} else {
		msg.Mode = "static"
	}
	msg.Address = gConf.Address
	msg.Netmask = gConf.Netmask
	msg.Gateway = gConf.Gateway
	msg.MainDNS = gConf.MainDNS
	msg.SecondaryDNS = gConf.SecondaryDNS

	setConfig(&msg)

	return nil
}

func (c *Config) Save() {
	config.Save()
}

var gConf Config

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	gConf.Dev = ""
	gConf.DHCP = true
	config.RegisterConfig("network", Init, Save)
}
