package agentgroup

import (
	"github.com/open-cmi/cmmns/module/agentgroup/router"
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthAPI("agentgroup", router.AuthGroup)
}
