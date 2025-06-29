package tools

import (
	"github.com/open-cmi/gobase/essential/webserver"
)

func init() {
	webserver.RegisterAuthRouter("tools", "/api/tools/v1")
}
