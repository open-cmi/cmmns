package captcha

import (
	"github.com/open-cmi/cmmns/module/rbac"
)

func init() {
	rbac.RegisterUnauthRouter("captcha", "/api/captcha/v1")
	rbac.UnauthAPI("captcha", "GET", "/", GetID)
	rbac.UnauthAPI("captcha", "GET", "/:id", GetPic)
}
