package secretkey

import (
	"github.com/open-cmi/cmmns/module/secretkey/router"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthAPI("secretkey", router.AuthGroup)
}
