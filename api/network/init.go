package network

import "github.com/open-cmi/cmmns/service/webserver"

func init() {
	webserver.RegisterAuthRouter("network", "/api/network/v1/")
	webserver.RegisterAuthAPI("network", "GET", "/", GetNetwork)
	webserver.RegisterAuthAPI("network", "PUT", "/", SetNetwork)
	webserver.RegisterAuthAPI("network", "GET", "/status/", GetNetworkStatus)
	webserver.RegisterAuthAPI("network", "POST", "/blinking/", BlinkingNetworkInterface)
}
