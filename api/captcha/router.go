package captcha

import (
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterUnauthRouter("captcha", "/api/captcha/v1")
	webserver.RegisterUnauthAPI("captcha", "GET", "/", GetID)
	webserver.RegisterUnauthAPI("captcha", "GET", "/:id", GetPic)
}
