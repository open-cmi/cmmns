package auditlog

import (
	"github.com/open-cmi/cmmns/module/rbac"
)

func init() {
	rbac.RegisterMustAuthRouter("auditlog", "/api/auditlog/v1")
}
