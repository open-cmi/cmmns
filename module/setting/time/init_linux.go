package time

import (
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/initial"
)

func Init() error {
	m := Get()
	if m == nil {
		switch gConf.Manage {
		case "systemd-timesyncd":
			var req SettingRequest
			req.AutoAdjust = true
			if len(gConf.NTPServers) != 0 {
				req.NtpServer = gConf.NTPServers[0] //"cn.ntp.org.cn"
			}
			req.TimeZone = "Asia/Shanghai"
			err := SetTimeSetting(&req)
			if err != nil {
				logger.Errorf("set time setting failed: %s\n", err.Error())
				return err
			}
		case "ntpd":
			var req SettingRequest
			req.AutoAdjust = true
			if len(gConf.NTPServers) != 0 {
				req.NtpServer = gConf.NTPServers[0] //"cn.ntp.org.cn"
			}
			req.TimeZone = "Asia/Shanghai"
			err := SetTimeSetting(&req)
			if err != nil {
				logger.Errorf("set time setting failed: %s\n", err.Error())
				return err
			}
		}
	}
	return nil
}

func init() {
	initial.Register("time", initial.PhaseDefault, Init)
}
