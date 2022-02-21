package def

type RegisterMsg struct {
	Token        string `json:"token"`
	LocalAddress string `json:"local_address"`
	Hostname     string `json:"hostname"`
	DevID        string `json:"dev_id"`
}
