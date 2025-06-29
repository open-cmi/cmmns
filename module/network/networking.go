package network

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/open-cmi/gobase/essential/logger"
)

type NetworkingConfig struct {
	Dev        string
	DHCP       bool
	Address    string
	Netmask    string
	Gateway    string
	NameServer []string
}

func NetworkingApply() error {
	// 写入文件
	var filename string
	if gConf.ConfFile == "" {
		filename = "/etc/network/interfaces"
	} else {
		filename = gConf.ConfFile
	}

	wf, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		logger.Errorf("open %s failed: %s\n", filename, err.Error())
		return err
	}

	var nicConfs []NetworkingConfig
	for _, dev := range gConf.Devices {
		devConf := GetNetConfig(dev)
		if devConf == nil {
			continue
		}
		if devConf.DHCP {
			var nwconf NetworkingConfig
			nwconf.Dev = dev
			nwconf.DHCP = true
			nicConfs = append(nicConfs, nwconf)
		} else {
			if devConf.Address == "" {
				continue
			}

			var nwconf NetworkingConfig
			nwconf.Dev = dev
			nwconf.DHCP = false
			nwconf.Address = devConf.Address
			nwconf.Netmask = devConf.Netmask

			if devConf.PreferredDNS != "" {
				nwconf.NameServer = append(nwconf.NameServer, devConf.PreferredDNS)
			}
			if devConf.AlternateDNS != "" {
				nwconf.NameServer = append(nwconf.NameServer, devConf.AlternateDNS)
			}
			nwconf.Gateway = devConf.Gateway
			nicConfs = append(nicConfs, nwconf)
		}
	}

	if len(nicConfs) != 0 {
		wf.WriteString("auto lo")
		wf.WriteString("\n")
		wf.WriteString("iface lo inet loopback")
		wf.WriteString("\n")

		for _, conf := range nicConfs {
			if conf.DHCP {
				msg := fmt.Sprintf("auto %s", conf.Dev)
				wf.WriteString(msg)
				wf.WriteString("\n")

				msg = fmt.Sprintf("iface %s inet dhcp", conf.Dev)
				wf.WriteString(msg)
				wf.WriteString("\n")
			} else if conf.Address != "" {
				msg := fmt.Sprintf("auto %s", conf.Dev)
				wf.WriteString(msg)
				wf.WriteString("\n")

				msg = fmt.Sprintf("iface %s inet static", conf.Dev)
				wf.WriteString(msg)
				wf.WriteString("\n")
				// address
				msg = fmt.Sprintf("address %s", conf.Address)
				wf.WriteString(msg)
				wf.WriteString("\n")
				// netmask
				msg = fmt.Sprintf("netmask %s", conf.Netmask)
				wf.WriteString(msg)
				wf.WriteString("\n")
				// gateway
				if conf.Gateway != "" {
					msg = fmt.Sprintf("gateway %s", conf.Gateway)
					wf.WriteString(msg)
					wf.WriteString("\n")
				}
				if len(conf.NameServer) != 0 {
					for _, ns := range conf.NameServer {
						msg = fmt.Sprintf("nameserver %s", ns)
						wf.WriteString(msg)
						wf.WriteString("\n")
					}
				}
			}
		}

		var args []string = []string{"restart", "networking"}
		err = exec.Command("systemctl", args...).Run()

		return err
	}
	return nil
}
