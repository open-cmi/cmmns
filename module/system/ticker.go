package system

import (
	"github.com/open-cmi/cmmns/service/ticker"
)

func init() {
	// ticker.Register("system_status", "0 */2 * * * *", func(name string, data interface{}) {
	// 	StartMonitor()
	// }, nil)
	ticker.Register("system_monitor", "*/15 * * * * *", func(name string, data interface{}) {
		MonitorCpuUsage()
		MonitorMemUsage()
	}, nil)
}
