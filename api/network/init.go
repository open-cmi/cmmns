package network

import "github.com/open-cmi/cmmns/service/webserver"

func init() {
	webserver.RegisterAuthRouter("system", "/api/common/v3/network/")
	webserver.RegisterAuthAPI("system", "GET", "/", GetNetwork)
	webserver.RegisterAuthAPI("system", "PUT", "/", SetNetwork)
}
