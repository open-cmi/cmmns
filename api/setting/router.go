package setting

import (
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("setting", "/api/common/v3/setting")
	webserver.RegisterAuthAPI("setting", "GET", "/", List)
	webserver.RegisterAuthAPI("setting", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("setting", "PUT", "/:id", Edit)
	webserver.RegisterAuthAPI("setting", "GET", "/pubnet/", GetPublicNet)
	webserver.RegisterAuthAPI("setting", "POST", "/pubnet/", SetPublicNet)
}
