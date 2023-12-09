package email

type SetRequest struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Sender   string `json:"sender"`
	Password string `json:"password"`
	UseTLS   bool   `json:"use_tls"`
}
