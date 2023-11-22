package system

import (
	"github.com/open-cmi/cmmns/api/system/setting"
	"github.com/open-cmi/cmmns/api/system/sysinfo"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("system-setting", "/api/system-setting/v1")

	webserver.RegisterAuthAPI("system-setting", "GET", "/device/", setting.GetDevID)
	webserver.RegisterAuthAPI("system-setting", "POST", "/reboot/", setting.Reboot)
	webserver.RegisterAuthAPI("system-setting", "POST", "/shutdown/", setting.ShutDown)
	webserver.RegisterAuthAPI("system-setting", "PUT", "/locale/", setting.ChangeLang)
	webserver.RegisterAuthAPI("system-setting", "GET", "/locale/", setting.GetLang)
	webserver.RegisterAuthAPI("system-setting", "POST", "/email/", setting.SetEmail)
	webserver.RegisterAuthAPI("system-setting", "GET", "/email/", setting.GetEmail)

	webserver.RegisterAuthRouter("system-info", "/api/system-info/v1")
	webserver.RegisterAuthAPI("system-info", "GET", "/basic-info/", sysinfo.GetBasicSystemInfo)
}
