package msg

// RegisterMsg register agent msg
type RegisterMsg struct {
	Token        string `json:"token"`
	DevID        string `json:"dev_id"`
	LocalAddress string `json:"local_address"`
	HostName     string `json:"hostname"`
}

// CreateMsg create agent msg
type CreateMsg struct {
	Address     string `json:"address"`
	Group       string `json:"group"`
	Port        int    `json:"port"`
	ConnType    string `json:"conn_type"`
	UserName    string `json:"username"`
	Passwd      string `json:"password"`
	SecretKey   string `json:"secret_key"`
	Description string `json:"description"`
}

type EditMsg struct {
	Address     string `json:"address"`
	Group       string `json:"group"`
	Port        int    `json:"port"`
	ConnType    string `json:"conn_type"`
	UserName    string `json:"username"`
	Passwd      string `json:"password"`
	SecretKey   string `json:"secret_key"`
	Description string `json:"description"`
}

// DeployMsg deploy msg
type DeployMsg struct {
	ID []string `json:"ids"`
}
