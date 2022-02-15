package user

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/storage/rdb"
	"github.com/open-cmi/cmmns/module/user/router"
)

func init() {

	api.RegisterAuthAPI("user", router.AuthGroup)
	api.RegisterUnauthAPI("user", router.UnauthGroup)
	rdb.Register("user", 3)
}
