package system

import (
	_ "github.com/open-cmi/cmmns/api/system/license"
	_ "github.com/open-cmi/cmmns/api/system/prod"
	_ "github.com/open-cmi/cmmns/api/system/setting"
	_ "github.com/open-cmi/cmmns/api/system/sysinfo"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("system", "/api/system/v1")
	webserver.RegisterUnauthRouter("system", "/api/system/v1")
}
