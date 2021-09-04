package config

import (
	"errors"

	"github.com/open-cmi/goutils/confparser"
)

// ConfParser conf parser
var ConfParser *confparser.Parser

// DatabaseModel database model
type DatabaseModel struct {
	Type     string `json:"type"`
	File     string `json:"file,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

// RedisCache redis cache
type RedisCache struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
}

// EmailInfo email info
type EmailInfo struct {
	From       string `json:"from"`
	SMTPServer string `json:"smtp_server"`
	User       string `json:"user"`
	Password   string `json:"password"`
	SMTPHost   string `json:"smtp_host"`
}

// Config config
type Config struct {
	Debug    bool          `json:"debug"`
	UnixSock bool          `json:"unixsock"`
	Model    DatabaseModel `json:"model"`
	Rdb      RedisCache    `json:"rdb"`
	Email    EmailInfo     `json:"email"`
	Domain   string        `json:"domain"`
}

var config Config

// Init config init
func Init(configfile string) (err error) {
	parser := confparser.New(configfile)
	if parser == nil {
		return errors.New("parse config failed")
	}
	err = parser.Load(&config)
	ConfParser = parser
	return err
}

// Save save config
func Save(c *Config) {
	ConfParser.Save(c)
}

// GetConfig get config
func GetConfig() *Config {
	return &config
}
