package agent

// CreateMsg create agent msg
type CreateMsg struct {
	Name        string `json:"name"`
	Group       int    `json:"servergroup"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	ConnType    string `json:"conntype"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
	SecretKey   string `json:"secretkey"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

// DeployMsg deploy msg
type DeployMsg struct {
	NodeID []string `json:"servers"`
}
