package network

import "github.com/open-cmi/cmmns/essential/config"

type Config struct {
	Dev          string `json:"dev"`
	DHCP         bool   `json:"dhcp"`
	Address      string `json:"address,omitempty"`
	Netmask      string `json:"netmask,omitempty"`
	Gateway      string `json:"gateway,omitempty"`
	MainDNS      string `json:"main_dns,omitempty"`
	SecondaryDNS string `json:"secondary_dns,omitempty"`
}

func (c *Config) Init() error {
	if c.Dev == "" {
		return nil
	}

	var msg ConfigMsg
	if c.DHCP {
		msg.Mode = "dhcp"
	} else {
		msg.Mode = "static"
	}
	msg.Address = c.Address
	msg.Netmask = c.Netmask
	msg.Gateway = c.Gateway
	msg.MainDNS = c.MainDNS
	msg.SecondaryDNS = c.SecondaryDNS

	setConfig(&msg)
	return nil
}

func (c *Config) Save() {
	config.Save()
}

var gConf Config

func init() {
	gConf.Dev = ""
	gConf.DHCP = true
	config.RegisterConfig("network", &gConf)
}
