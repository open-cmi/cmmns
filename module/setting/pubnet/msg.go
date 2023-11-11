package pubnet

type SetPublicNetRequest struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Schema string `json:"schema"`
}
