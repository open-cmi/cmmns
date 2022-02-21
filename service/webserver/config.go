package webserver

import (
	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/service/webserver/middleware"
)

type Config struct {
	Middleware middleware.MiddlewareConfig `json:"middleware"`
	Listen     string                      `json:"listen"`
	Port       int                         `json:"port"`
	UnixPath   string                      `json:"unix_path"`
}

var moduleConfig Config

func (c *Config) Init() error {
	return nil
}

func init() {
	// default config
	moduleConfig.Listen = "127.0.0.1"
	moduleConfig.Port = 30000
	moduleConfig.UnixPath = "/tmp/cmmns.sock"
	moduleConfig.Middleware.SessionStore = "memory"

	config.RegisterConfig("webserver", &moduleConfig)
}
