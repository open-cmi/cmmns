package secretkey

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/secretkey/router"
)

func init() {
	api.RegisterAuthAPI("secretkey", router.AuthGroup)
}
