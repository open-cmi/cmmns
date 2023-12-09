package assist

import (
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("assist", "/api/assist/v1")
	webserver.RegisterAuthAPI("assist", "GET", "/", GetAssist)
	webserver.RegisterAuthAPI("assist", "POST", "/", SetAssist)
}
