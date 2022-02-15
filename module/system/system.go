package system

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/ticker"
	"github.com/open-cmi/cmmns/module/system/model"
	"github.com/open-cmi/cmmns/module/system/router"
)

func init() {
	ticker.RegisterTicker("system_status", "0 */5 * * * *", func() {
		model.StartMonitor()
	})
	api.RegisterAuthAPI("system", router.AuthGroup)
}
