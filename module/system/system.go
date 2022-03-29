package system

import (
	"github.com/open-cmi/cmmns/service/ticker"
)

func init() {
	ticker.Register("system_status", "0 */5 * * * *", func() {
		StartMonitor()
	})
}
