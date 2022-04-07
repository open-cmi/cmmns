package user

import "github.com/open-cmi/cmmns/essential/config"

// Config smtp
type Config struct {
	From       string `json:"from"`
	SMTPServer string `json:"smtp_server"`
	User       string `json:"user"`
	Password   string `json:"password"`
	SMTPHost   string `json:"smtp_host"`
	Domain     string `json:"domain"`
}

func (c *Config) Init() error {
	return nil
}

var gConf Config

func init() {
	config.RegisterConfig("smtp", &gConf)
}
