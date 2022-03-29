package agent

type Config struct {
	LinuxPackage string `json:"linux_package"` // linux packge的位置，相对root path路径
	LocalAddress string `json:"local_address"`
	LocalPort    int    `json:"local_port"`
	LocalProto   string `json:"local_proto"`
	Address      string `json:"address"`
	Port         int    `json:"port"`
	Proto        string `json:"proto"`
}

func (c *Config) Init() error {
	return nil
}
