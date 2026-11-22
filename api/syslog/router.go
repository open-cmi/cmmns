package syslog

import (
	"github.com/open-cmi/cmmns/module/rbac"
)

func init() {
	rbac.RegisterMustAuthRouter("syslog", "/api/syslog/v1")
}
