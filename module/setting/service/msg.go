package service

type SetServicePortRequest struct {
	HTTPPort  int `json:"http_port"`
	HTTPSPort int `json:"https_port"`
}
