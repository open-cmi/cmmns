package webserver

import (
	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/service/webserver/middleware"
)

type Server struct {
	Address  string `json:"address"`
	Port     int    `json:"port,omitempty"`
	Proto    string `json:"proto"`
	CertFile string `json:"cert,omitempty"`
	KeyFile  string `json:"key,omitempty"`
}

type Config struct {
	Middleware middleware.MiddlewareConfig `json:"middleware"`
	Server     []Server                    `json:"server"`
}

var moduleConfig Config

func (c *Config) Init() error {
	return nil
}

func init() {
	// default config
	moduleConfig.Middleware.SessionStore = "memory"

	moduleConfig.Server = append(moduleConfig.Server, Server{
		Address: "127.0.0.1",
		Port:    30000,
		Proto:   "http",
	}, Server{
		Address: "/tmp/cmmns.sock",
		Proto:   "unix",
	})

	config.RegisterConfig("webserver", &moduleConfig)
}
