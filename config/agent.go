package config

// AgentConfig  agent config struct
type AgentConfig struct {
	Debug       bool   `json:"debug"`
	UnixSock    bool   `json:"unixsock"`
	Master      string `json:"master"`
	APIServer   string `json:"api_server"`
	V2RayConfig string `json:"v2ray_config"`
}
