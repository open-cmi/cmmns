package network

import (
	"net"
	"sort"

	"github.com/open-cmi/cmmns/essential/config"
)

func Get() []ConfigRequest {
	var mode string = "static"
	var devices []ConfigRequest
	for dev := range gConf.Devices {
		var conf ConfigRequest
		v := gConf.Devices[dev]
		if v.DHCP {
			mode = "dhcp"
		}

		conf.Address = v.Address
		conf.Netmask = v.Netmask
		conf.Dev = dev
		conf.Gateway = v.Gateway
		conf.Mode = mode
		conf.PreferredDNS = v.PreferredDNS
		conf.AlternateDNS = v.AlternateDNS
		devices = append(devices, conf)
	}

	sort.SliceStable(devices, func(i int, j int) bool {
		dev1 := devices[i]
		dev2 := devices[j]
		return dev1.Dev < dev2.Dev
	})
	return devices
}

func GetInterfaceStatus(f net.Flags) string {
	s := ""
	if f&net.FlagUp != 0 {
		if s != "" {
			s += "|"
		}
		s += "up"
	}
	if f&net.FlagRunning != 0 {
		if s != "" {
			s += "|"
		}
		s += "running"
	}

	if s == "" {
		s = "0"
	}
	return s
}

func GetStatus() ([]InterfaceStatus, error) {
	var resp []InterfaceStatus
	for dev := range gConf.Devices {
		var status InterfaceStatus
		status.Dev = dev

		intf, err := net.InterfaceByName(dev)
		if err != nil {
			return resp, err
		}
		addrs, err := intf.Addrs()
		if err != nil {
			return resp, err
		}
		status.MTU = intf.MTU
		status.EtherAddr = intf.HardwareAddr.String()
		status.Status = GetInterfaceStatus(intf.Flags)
		if len(addrs) != 0 {
			var ipv4 string
			var masklen int
			for _, addr := range addrs {
				ip, ipnet, _ := net.ParseCIDR(addr.String())
				if ip.To4() != nil {
					ipv4 = ip.To4().String()
					masklen, _ = ipnet.Mask.Size()
				}
			}
			netmask := NetmaskString(masklen)
			status.Address = ipv4
			status.Netmask = netmask
		}
		// 取网卡的统计信息
	}

	return resp, nil
}

func GetRoutes() {
	// gw := GetDefaultGateway()

	// dns := GetDNS()
	// var mainDNS string
	// var secondaryDNS string
	// if len(dns) >= 1 {
	// 	mainDNS = dns[0]
	// }
	// if len(dns) >= 2 {
	// 	secondaryDNS = dns[1]
	// }
	// return ConfigRequest{
	// 	Address:      ipv4,
	// 	Netmask:      netmask,
	// 	Gateway:      gw.String(),
	// 	PreferredDNS: mainDNS,
	// 	AlternateDNS: secondaryDNS,
	// }
}

func Set(msg *ConfigRequest) error {
	for name := range gConf.Devices {
		// 接口相同，或者dev为空，取第一个
		if name == msg.Dev {
			conf := gConf.Devices[name]
			if msg.Mode == "dhcp" {
				conf.DHCP = true
				conf.Address = ""
				conf.Netmask = ""
				conf.Gateway = ""
				conf.PreferredDNS = ""
				conf.AlternateDNS = ""
			} else {
				conf.DHCP = false
				conf.Address = msg.Address
				conf.Netmask = msg.Netmask
				conf.Gateway = msg.Gateway
				conf.PreferredDNS = msg.PreferredDNS
				conf.AlternateDNS = msg.AlternateDNS
			}

			gConf.Devices[name] = conf

			err := NetworkApply(&gConf)
			if err != nil {
				return err
			}
			config.Save()
			break
		}
	}

	return nil
}
