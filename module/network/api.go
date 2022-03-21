package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/open-cmi/cmmns/essential/logger"
	"gopkg.in/yaml.v2"
)

func Get() ConfigMsg {
	var mode string = "static"
	if gConf.DHCP {
		mode = "dhcp"
	}

	return ConfigMsg{
		Mode:         mode,
		Address:      gConf.Address,
		Netmask:      gConf.Netmask,
		Gateway:      gConf.Gateway,
		MainDNS:      gConf.MainDNS,
		SecondaryDNS: gConf.SecondaryDNS,
	}
}

func GetStatus() ConfigMsg {
	var conf ConfigMsg
	intf, err := net.InterfaceByName("en0")
	if err != nil {
		return conf
	}
	addrs, err := intf.Addrs()
	if err != nil || len(addrs) == 0 {
		return conf
	}
	fmt.Println(addrs)
	var ipv4 string
	var masklen int
	for _, addr := range addrs {
		ip, ipnet, _ := net.ParseCIDR(addr.String())
		if ip.To4() != nil {
			ipv4 = ip.To4().String()
			masklen, _ = ipnet.Mask.Size()
		}
	}
	gw := GetDefaultGateway()
	netmask := NetmaskString(masklen)
	dns := GetDNS()
	var mainDNS string
	var secondaryDNS string
	if len(dns) >= 1 {
		mainDNS = dns[0]
	}
	if len(dns) >= 2 {
		secondaryDNS = dns[1]
	}
	return ConfigMsg{
		Address:      ipv4,
		Netmask:      netmask,
		Gateway:      gw.String(),
		MainDNS:      mainDNS,
		SecondaryDNS: secondaryDNS,
	}
}

func Set(msg *ConfigMsg) error {
	err := setConfig(msg)
	if err != nil {
		return err
	}

	gConf.Address = msg.Address
	gConf.Netmask = msg.Netmask
	gConf.Gateway = msg.Gateway
	gConf.MainDNS = msg.MainDNS
	gConf.SecondaryDNS = msg.SecondaryDNS
	if msg.Mode == "dhcp" {
		gConf.DHCP = true
	} else {
		gConf.DHCP = false
	}

	gConf.Save()
	return nil
}

func setConfig(msg *ConfigMsg) error {
	// 这里要校验格式

	// 写入文件
	filename := fmt.Sprintf("/tmp/99-%s.yaml", gConf.Dev)
	wf, err := os.OpenFile(filename,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open netplan config failed: %s\n", err.Error())
		return err
	}
	var netconf NetConfig
	netconf.Network.Version = 2
	netconf.Network.Renderer = "networkd"
	netconf.Network.Ethernets = make(map[string]EthernetConfig)
	if msg.Mode == "dhcp" {
		netconf.Network.Ethernets[gConf.Dev] = EthernetConfig{
			DHCP4: "yes",
		}
	} else {
		mask := net.IPMask(net.ParseIP(msg.Netmask).To4()) // If you have the mask as a string
		maskLen, _ := mask.Size()
		addr := fmt.Sprintf("%s/%d", msg.Address, maskLen)
		dns := []string{}
		if msg.MainDNS != "" {
			dns = append(dns, msg.MainDNS)
		}
		if msg.SecondaryDNS != "" {
			dns = append(dns, msg.SecondaryDNS)
		}
		netconf.Network.Ethernets[gConf.Dev] = EthernetConfig{
			DHCP4: "no",
			Addresses: []string{
				addr,
			},
			Gateway4:    msg.Gateway,
			Nameservers: dns,
		}
	}

	netout, err := yaml.Marshal(netconf)
	if err != nil {
		return err
	}
	wf.Write(netout)

	var args []string = []string{"apply"}
	exec.Command("netplan", args...).Run()

	return nil
}
