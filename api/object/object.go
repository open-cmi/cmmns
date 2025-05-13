package object

import (
	_ "github.com/open-cmi/cmmns/api/object/time"
	"github.com/open-cmi/cmmns/essential/webserver"
)

func init() {
	webserver.RegisterAuthRouter("object", "/api/object/v1")
}
