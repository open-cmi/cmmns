package tools

import "github.com/open-cmi/cmmns/module/rbac"

func init() {
	rbac.RegisterOptionAuthRouter("tools", "/api/tools/v1")
}
