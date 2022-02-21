package auditlog

import (
	"github.com/open-cmi/cmmns/module/auditlog/router"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthAPI("auditlog", router.AuthGroup)
}
