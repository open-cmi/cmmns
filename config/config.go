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

// Config config
type Config struct {
	Debug    bool          `json:"debug"`
	UnixSock bool          `json:"unixsock"`
	Model    DatabaseModel `json:"model"`
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
