package auditlog

import (
	"github.com/open-cmi/cmmns/essential/webserver"
)

func init() {
	webserver.RegisterMustAuthRouter("auditlog", "/api/auditlog/v1")
	webserver.RegisterMustAuthAPI("auditlog", "GET", "/", List)
}
