package portal

import (
	"github.com/open-cmi/cmmns/module/rbac"
)

func init() {
	rbac.RegisterMustAuthRouter("portal", "/api/portal/v1")
}
