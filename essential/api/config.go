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
	config.RegisterConfig("api", &moduleConfig)
}
