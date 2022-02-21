package agent

import (
	"github.com/open-cmi/cmmns/module/agent/router"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {

	webserver.RegisterAuthAPI("agent", router.AuthGroup)
	webserver.RegisterUnauthAPI("agent", router.UnauthGroup)
}
