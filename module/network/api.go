package network

import (
	"errors"
	"net"
	"os/exec"
	"sort"

	netio "github.com/shirou/gopsutil/net"

	"github.com/open-cmi/cmmns/essential/logger"
)

func Get() []ConfigRequest {
	var mode string = "static"
	var devices []ConfigRequest
	for _, dev := range gConf.Devices {
		var conf ConfigRequest
		v := GetNetConfig(dev)
		if v == nil {
			conf.Mode = "static"
			conf.Dev = dev
		} else {
			if v.DHCP {
				mode = "dhcp"
			}
			conf.Dev = dev
			conf.Mode = mode

			conf.Address = v.Address
			conf.Netmask = v.Netmask
			conf.Gateway = v.Gateway
			conf.PreferredDNS = v.PreferredDNS
			conf.AlternateDNS = v.AlternateDNS
		}

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
		s += "up"
	} else {
		s += "down"
	}
	if f&net.FlagRunning != 0 {
		if s != "" {
			s += "|"
		}
		s += "running"
	}

	return s
}

func GetStatus() (int, []InterfaceStatus, error) {
	var resp []InterfaceStatus
	counters, err := netio.IOCounters(true)
	if err != nil {
		return 0, resp, err
	}

	for _, dev := range gConf.Devices {
		var status InterfaceStatus
		status.Dev = dev

		intf, err := net.InterfaceByName(dev)
		if err != nil {
			return 0, resp, err
		}
		addrs, err := intf.Addrs()
		if err != nil {
			return 0, resp, err
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
		for _, counter := range counters {
			if counter.Name == dev {
				status.BytesRecv = counter.BytesRecv
				status.BytesSent = counter.BytesSent
				status.PacketsRecv = counter.PacketsRecv
				status.PacketsSent = counter.PacketsSent
			}
		}

		resp = append(resp, status)
	}

	return len(resp), resp, nil
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
	var found bool = false
	for _, name := range gConf.Devices {
		// 接口相同，或者dev为空，取第一个
		if name == msg.Dev {
			found = true
			break
		}
	}
	if !found {
		return errors.New("dev not supported")
	}

	m := GetNetConfig(msg.Dev)
	if m == nil {
		m = New()
	}
	m.Dev = msg.Dev
	if msg.Mode == "dhcp" {
		m.DHCP = true
		m.Address = ""
		m.Netmask = ""
		m.Gateway = ""
		m.PreferredDNS = ""
		m.AlternateDNS = ""
	} else {
		m.DHCP = false
		m.Address = msg.Address
		m.Netmask = msg.Netmask
		m.Gateway = msg.Gateway
		m.PreferredDNS = msg.PreferredDNS
		m.AlternateDNS = msg.AlternateDNS
	}
	err := m.Save()
	if err != nil {
		return err
	}

	err = NetworkApply(&gConf)
	return err
}

func BlinkingInterface(req *BlinkingRequest) error {
	args := []string{"-p", req.Dev, "5"}
	cmd := exec.Command("ethtool", args...)
	err := cmd.Run()
	if err != nil {
		logger.Errorf("run ethtool failed: %s\n", err.Error())
		return errors.New("operation not supported")
	}
	return nil
}

type SetManagementInterfaceRequest struct {
	Devices []string `json:"devices"`
}

func SetManagementInterface(req *SetManagementInterfaceRequest) error {
	m := GetManagementInterfaceModel()
	if m == nil {
		m = &ManagementInterfaceModel{
			isNew: true,
		}
	}
	m.Interfaces = req.Devices
	err := m.Save()
	if err != nil {
		return nil
	}
	gConf.Devices = req.Devices
	return nil
}

func GetAvailableManagementInterface() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return []string{}, err
	}
	var devices []string = []string{}
	for _, intf := range interfaces {
		if intf.Name == "lo" {
			continue
		}
		devices = append(devices, intf.Name)
	}
	return devices, nil
}

func IsManagementInterface(dev string) bool {
	for _, device := range gConf.Devices {
		if device == dev {
			return true
		}
	}
	return false
}
