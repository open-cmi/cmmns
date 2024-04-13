package network

import (
	"testing"
)

func TestSet(t *testing.T) {

	var msg ConfigRequest
	msg.Mode = "dhcp"
	err := Set(&msg)
	if err != nil {
		t.Errorf("set net config failed")
	}

	msg.Mode = "static"
	msg.Address = "192.168.56.2"
	msg.Netmask = "255.0.0.0"
	msg.Gateway = "192.168.56.1"
	msg.PreferredDNS = "8.8.8.8"
	err = Set(&msg)
	if err != nil {
		t.Errorf("set net config failed")
	}
}
