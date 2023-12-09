package system

import (
	_ "github.com/open-cmi/cmmns/api/system/setting"
	_ "github.com/open-cmi/cmmns/api/system/sysinfo"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("system-setting", "/api/system-setting/v1")
	webserver.RegisterAuthRouter("system-info", "/api/system-info/v1")
}
