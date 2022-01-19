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
	AgentPackageLocation string `json:"agent_package_location"` // agent packge的本地位置
	AgentConfigLocation  string `json:"agent_config_location"`  // agent的远端配置文件位置
}

type MasterInfoConfig struct {
	InternalAddress string `json:"internal_address"`
	InternalPort    int    `json:"internal_port"`
	InternalProto   string `json:"internal_proto"`
	ExternalAddress string `json:"external_address"`
	ExternalPort    int    `json:"external_port"`
	ExternalProto   string `json:"external_proto"`
}

type RemoteService struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	LocalIP    string `json:"local_ip"`
	LocalPort  uint16 `json:"local_port"`
	RemotePort uint16 `json:"remote_port"`
}

type RemoteAssistConfig struct {
	ServerAddr string          `json:"server_addr"`
	ServerPort uint16          `json:"server_port"`
	Token      string          `json:"token,omitempty"`
	Service    []RemoteService `json:"services"`
}

// Config config
type Config struct {
	Debug        bool               `json:"debug"`
	UnixSock     bool               `json:"unixsock"`
	StoreType    string             `json:"storetype,omitempty"`
	Model        DatabaseModel      `json:"model"`
	Rdb          RedisCache         `json:"rdb"`
	Email        EmailInfo          `json:"email"`
	MasterInfo   MasterInfoConfig   `json:"masterinfo"`
	Distributed  DistributedConfig  `json:"distributed"`
	RemoteAssist RemoteAssistConfig `json:"remote_assist"`
}

var config Config

// Init config init
func Init(configfile string) (err error) {
	parser := confparser.New(configfile)
	if parser == nil {
		return errors.New("parse config failed")
	}
	err = parser.Load(&config)
	if config.StoreType == "" {
		config.StoreType = "redis"
	}
	ConfParser = parser
	return err
}

// Save save config
func (c *Config) Save() {
	ConfParser.Save(c)
}

// GetConfig get config
func GetConfig() *Config {
	return &config
}
