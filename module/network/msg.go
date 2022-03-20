package network

type ConfigMsg struct {
	Mode         string `json:"mode"`
	Address      string `json:"address"`
	Netmask      string `json:"netmask"`
	Gateway      string `json:"gateway"`
	MainDNS      string `json:"main_dns"`
	SecondaryDNS string `json:"secondary_dns"`
}
