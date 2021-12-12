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

type DistributedConfig struct {
	InternalAddress string `json:"internal_address"`
	InternalPort    int    `json:"internal_port"`
	InternalProto   string `json:"internal_proto"`
	ExternalAddress string `json:"external_address"`
	ExternalPort    int    `json:"external_port"`
	ExternalProto   string `json:"external_proto"`
	AgentLocation   string `json:"agent_location"`
}

// Config config
type Config struct {
	Debug            bool              `json:"debug"`
	UnixSock         bool              `json:"unixsock"`
	Model            DatabaseModel     `json:"model"`
	Rdb              RedisCache        `json:"rdb"`
	Email            EmailInfo         `json:"email"`
	Distributed      DistributedConfig `json:"distributed"`
	MasterInfoConfig string            `json:"master_info_config"`
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
