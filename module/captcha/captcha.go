package captcha

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/captcha/router"
)

func init() {
	api.RegisterUnauthAPI("captcha", router.UnauthGroup)
}
