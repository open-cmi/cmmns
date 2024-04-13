package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"

	"github.com/open-cmi/cmmns/essential/logger"
	"gopkg.in/yaml.v2"
)

var NetPlanDir string = "/etc/netplan/"

type NameServerConfig struct {
	Addresses []string `json:"addresses,omitempty"`
}

type EthernetConfig struct {
	DHCP4       string           `yaml:"dhcp4"`
	Addresses   []string         `yaml:"addresses,omitempty"`
	Gateway4    string           `yaml:"gateway4,omitempty"`
	Nameservers NameServerConfig `yaml:"nameservers,omitempty"`
	Interfaces  []string         `yaml:"interfaces,omitempty"`
}

type Network struct {
	Version   int                       `yaml:"version"`
	Renderer  string                    `yaml:"renderer"`
	Ethernets map[string]EthernetConfig `yaml:"ethernets"`
	Bridges   string                    `yaml:"bridges,omitempty"`
}

type NetplanConfig struct {
	Network Network `yaml:"network"`
}

func NetplanApply(conf *Config) error {
	// 写入文件
	var filename string
	if conf.ConfFile == "" {
		files, err := os.ReadDir(NetPlanDir)
		if err != nil {
			return err
		}
		// 取第一个文件作为配置文件
		for _, file := range files {
			if !file.IsDir() {
				filename = file.Name()
				break
			}
		}
		// 如果无任何文件，则命名一个文件
		if filename == "" {
			filename = "./50-cloud-init.yaml"
		}
		filename = path.Join(NetPlanDir, filename)
	} else {
		filename = conf.ConfFile
	}
	wf, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open netplan config failed: %s\n", err.Error())
		return err
	}
	var netconf NetplanConfig
	netconf.Network.Version = 2
	netconf.Network.Renderer = "networkd"
	netconf.Network.Ethernets = make(map[string]EthernetConfig)
	for dev := range conf.Devices {
		devConf := conf.Devices[dev]
		if devConf.DHCP {
			netconf.Network.Ethernets[dev] = EthernetConfig{
				DHCP4: "yes",
			}
		} else {
			if devConf.Address == "" {
				continue
			}
			mask := net.IPMask(net.ParseIP(devConf.Netmask).To4()) // If you have the mask as a string
			maskLen, _ := mask.Size()
			addr := fmt.Sprintf("%s/%d", devConf.Address, maskLen)
			var dns NameServerConfig
			if devConf.PreferredDNS != "" {
				dns.Addresses = append(dns.Addresses, devConf.PreferredDNS)
			}
			if devConf.AlternateDNS != "" {
				dns.Addresses = append(dns.Addresses, devConf.AlternateDNS)
			}
			netconf.Network.Ethernets[dev] = EthernetConfig{
				DHCP4: "no",
				Addresses: []string{
					addr,
				},
				Gateway4:    devConf.Gateway,
				Nameservers: dns,
			}
		}
	}

	if len(netconf.Network.Ethernets) != 0 {
		netout, err := yaml.Marshal(netconf)
		if err != nil {
			return err
		}
		wf.Write(netout)

		var args []string = []string{"apply"}
		err = exec.Command("netplan", args...).Run()

		return err
	}
	return nil
}
