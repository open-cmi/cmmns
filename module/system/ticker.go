package system

import (
	"github.com/open-cmi/gobase/essential/ticker"
)

func MonitorSystemUsage() {
	if gConf.MonitorUsage {
		MonitorCpuUsage()
		MonitorMemUsage()
	}
}

func init() {
	// ticker.Register("system_status", "0 */2 * * * *", func(name string, data interface{}) {
	// 	StartMonitor()
	// }, nil)
	ticker.Register("system_monitor", "*/30 * * * * *", func(name string, data interface{}) {
		MonitorSystemUsage()
	}, nil)
}
