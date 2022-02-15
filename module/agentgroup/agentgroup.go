package agentgroup

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/agentgroup/router"
)

func init() {
	api.RegisterAuthAPI("agentgroup", router.AuthGroup)
}
