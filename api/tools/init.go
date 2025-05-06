package tools

import (
	"github.com/open-cmi/cmmns/essential/webserver"
)

func init() {
	webserver.RegisterAuthRouter("tools", "/api/tools/v1")
}
