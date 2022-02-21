package user

import (
	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/essential/storage/rdb"
	"github.com/open-cmi/cmmns/module/user/router"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {

	webserver.RegisterAuthAPI("user", router.AuthGroup)
	webserver.RegisterUnauthAPI("user", router.UnauthGroup)
	rdb.Register("user", def.RDBUser)
}
