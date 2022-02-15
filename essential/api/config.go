package api

import (
	"github.com/open-cmi/cmmns/essential/api/middleware"
	"github.com/open-cmi/cmmns/essential/config"
)

type Config struct {
	Middleware middleware.MiddlewareConfig `json:"middleware"`
	Listen     string                      `json:"listen"`
	Port       int                         `json:"port"`
	UnixPath   string                      `json:"unix_path"`
}

var moduleConfig Config

func init() {
	// default config
	moduleConfig.Listen = "127.0.0.1"
	moduleConfig.Port = 30000
	moduleConfig.UnixPath = "/tmp/cmmns.sock"
	moduleConfig.Middleware.SessionStore = "memory"

	config.RegisterConfig("api", &moduleConfig)
}
