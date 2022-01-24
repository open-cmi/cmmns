package agent

// CreateMsg create agent msg
type CreateMsg struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	ConnType    string `json:"conn_type"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
	SecretKey   string `json:"secret_key"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

type EditMsg struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	ConnType    string `json:"conn_type"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
	SecretKey   string `json:"secret_key"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

// DeployMsg deploy msg
type DeployMsg struct {
	ID []string `json:"ids"`
}
