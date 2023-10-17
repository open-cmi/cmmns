package system

import (
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("system", "/api/common/v3/system/")
	webserver.RegisterAuthAPI("system", "GET", "/status/", GetStatus)
	webserver.RegisterAuthAPI("system", "GET", "/device/", GetDevID)
	webserver.RegisterAuthAPI("system", "POST", "/reboot/", Reboot)
	webserver.RegisterAuthAPI("system", "POST", "/shutdown/", ShutDown)
	webserver.RegisterAuthAPI("system", "PUT", "/locale", ChangeLang)
	webserver.RegisterAuthAPI("system", "GET", "/locale", GetLang)
}
