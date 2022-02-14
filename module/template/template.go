package template

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/template/router"
)

func init() {
	api.RegisterAuthAPI("template", router.AuthGroup)
}
