package network

import (
	"bytes"
	"fmt"
	"strconv"
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
