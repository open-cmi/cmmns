package dnsmasq

import (
	"github.com/open-cmi/cmmns/pkg/shell"
	"github.com/open-cmi/gobase/essential/logger"
)

func RestartService() error {
	//shell.Execute("systemctl stop systemd-resolved")
	//defer shell.Execute("systemctl start systemd-resolved")
	err := shell.Execute(gConf.Restart)
	if err != nil {
		logger.Error("restart dnsmasq service failed\n")
		return err
	}

	return nil
}

func StartService() error {
	//shell.Execute("systemctl stop systemd-resolved")
	//defer shell.Execute("systemctl start systemd-resolved")
	err := shell.Execute(gConf.Start)
	if err != nil {
		logger.Error("start dnsmasq service failed\n")
		return err
	}

	return nil
}

func StopService() error {
	//shell.Execute("systemctl stop systemd-resolved")
	//defer shell.Execute("systemctl start systemd-resolved")
	err := shell.Execute(gConf.Stop)
	if err != nil {
		logger.Error("start dnsmasq service failed\n")
		return err
	}

	return nil
}
