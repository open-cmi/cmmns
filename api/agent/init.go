package agent

import (
	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/service/webserver"
)

var moduleConfig Config

func init() {
	config.RegisterConfig("cluster", &moduleConfig)
	webserver.RegisterAuthAPI("agent", AuthGroup)
	webserver.RegisterUnauthAPI("agent", UnauthGroup)
}
