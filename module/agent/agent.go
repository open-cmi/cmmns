package agent

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/agent/router"
)

func init() {

	api.RegisterAuthAPI("agent", router.AuthGroup)
	api.RegisterUnauthAPI("agent", router.UnauthGroup)
}
