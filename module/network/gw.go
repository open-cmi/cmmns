package network

import (
	"net"

	"github.com/jackpal/gateway"
)

func GetDefaultGateway() net.IP {
	gw, err := gateway.DiscoverGateway()
	if err != nil {
		return net.IPv4(0, 0, 0, 0)
	}
	return gw
}
