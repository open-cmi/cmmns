package captcha

import (
	"github.com/open-cmi/cmmns/module/captcha/router"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterUnauthAPI("captcha", router.UnauthGroup)
}
