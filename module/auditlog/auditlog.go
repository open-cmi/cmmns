package auditlog

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/auditlog/router"
)

func init() {
	api.RegisterAuthAPI("agent", router.AuthGroup)
}
