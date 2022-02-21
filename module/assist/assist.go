package assist

import (
	"github.com/open-cmi/cmmns/module/assist/router"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthAPI("assist", router.AuthGroup)
}
