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

type AgentConfig struct {
	LinuxPackage string `json:"linux_package"` // linux packge的位置，相对root path路径
}

type MasterConfig struct {
	LocalAddress string `json:"local_address"`
	LocalPort    int    `json:"local_port"`
	LocalProto   string `json:"local_proto"`
	Address      string `json:"address"`
	Port         int    `json:"port"`
	Proto        string `json:"proto"`
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
	Master       MasterConfig       `json:"master,omitempty"`
	Agent        AgentConfig        `json:"agent,omitempty"`
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
