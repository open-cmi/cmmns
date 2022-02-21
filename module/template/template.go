package template

import (
	"github.com/open-cmi/cmmns/module/template/router"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthAPI("template", router.AuthGroup)
}
