package system

import (
	"github.com/open-cmi/cmmns/module/system/model"
	"github.com/open-cmi/cmmns/module/system/router"
	"github.com/open-cmi/cmmns/service/ticker"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	ticker.Register("system_status", "0 */5 * * * *", func() {
		model.StartMonitor()
	})
	webserver.RegisterAuthAPI("system", router.AuthGroup)
}
