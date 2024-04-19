package network

type ConfigRequest struct {
	Dev          string `json:"dev"`
	Mode         string `json:"mode"`
	Address      string `json:"address"`
	Netmask      string `json:"netmask"`
	Gateway      string `json:"gateway"`
	PreferredDNS string `json:"preferred_dns"`
	AlternateDNS string `json:"alternate_dns"`
}

type InterfaceStatus struct {
	Dev         string `json:"dev"`
	Address     string `json:"address"`
	Netmask     string `json:"netmask"`
	Status      string `json:"status"`
	MTU         int    `json:"mtu"`
	EtherAddr   string `json:"ether_addr"`
	PacketsRecv uint64 `json:"packets_recv"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	BytesSent   uint64 `json:"bytes_sent"`
}

type BlinkingRequest struct {
	Dev string `json:"dev"`
}
