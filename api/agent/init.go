package agent

import (
	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/service/webserver"
)

var gConf Config

func init() {
	config.RegisterConfig("cluster", &gConf)
	webserver.RegisterAuthAPI("agent", AuthGroup)
	webserver.RegisterUnauthAPI("agent", UnauthGroup)
}
