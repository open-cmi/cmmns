package network

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/route"
)

var Location string

type EthernetConfig struct {
	DHCP4       string   `yaml:"dhcp4"`
	Addresses   []string `yaml:"addresses,omitempty"`
	Gateway4    string   `yaml:"gateway4,omitempty"`
	Nameservers []string `yaml:"nameservers,omitempty"`
	Interfaces  []string `yaml:"interfaces,omitempty"`
}

type Network struct {
	Version   int                       `yaml:"version"`
	Renderer  string                    `yaml:"renderer"`
	Ethernets map[string]EthernetConfig `yaml:"ethernets"`
	Bridges   string                    `yaml:"bridges,omitempty"`
}

type NetConfig struct {
	Network Network `yaml:"network"`
}

// 如 24 对应的子网掩码地址为 255.255.255.0
func NetmaskString(subnet int) string {
	var buff bytes.Buffer
	for i := 0; i < subnet; i++ {
		buff.WriteString("1")
	}
	for i := subnet; i < 32; i++ {
		buff.WriteString("0")
	}
	masker := buff.String()
	a, _ := strconv.ParseUint(masker[:8], 2, 64)
	b, _ := strconv.ParseUint(masker[8:16], 2, 64)
	c, _ := strconv.ParseUint(masker[16:24], 2, 64)
	d, _ := strconv.ParseUint(masker[24:32], 2, 64)
	resultMask := fmt.Sprintf("%v.%v.%v.%v", a, b, c, d)
	return resultMask
}

var defaultRoute = [4]byte{0, 0, 0, 0}

func GetDefaultGateway() (gw net.IP) {
	rib, _ := route.FetchRIB(0, route.RIBTypeRoute, 0)
	messages, err := route.ParseRIB(route.RIBTypeRoute, rib)

	if err != nil {
		return gw
	}

	var destination, gateway *route.Inet4Addr

	for _, message := range messages {
		route_message := message.(*route.RouteMessage)
		addresses := route_message.Addrs

		ok := false

		if destination, ok = addresses[0].(*route.Inet4Addr); !ok {
			continue
		}

		if gateway, ok = addresses[1].(*route.Inet4Addr); !ok {
			continue
		}

		if destination == nil || gateway == nil {
			continue
		}

		if destination.IP == defaultRoute {
			break
		}
	}
	if gateway != nil {
		return net.IPv4(gateway.IP[0], gateway.IP[1], gateway.IP[2], gateway.IP[3])
	}
	return net.IPv4(0, 0, 0, 0)
}

func GetDNS() []string {
	var dns []string
	rf, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return dns
	}
	rd := bufio.NewReader(rf)
	for linebyte, _, err := rd.ReadLine(); err == nil; linebyte, _, err = rd.ReadLine() {
		line := string(linebyte)
		fmt.Println(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "nameserver") {
			tmp := strings.TrimPrefix(line, "nameserver")
			tmp = strings.Trim(tmp, "\r\n\t ")
			dns = append(dns, tmp)
		}
	}
	return dns
}
