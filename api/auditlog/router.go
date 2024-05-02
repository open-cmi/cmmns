package auditlog

import (
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("auditlog", "/api/auditlog/v1")
	webserver.RegisterAuthAPI("auditlog", "GET", "/", List)
}
