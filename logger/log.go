package logger

import (
	"path/filepath"

	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/goutils/logutil"
	"github.com/open-cmi/goutils/pathutil"
)

const (
	Debug = iota
	Info
	Warn
	Error
)

var Logger *logutil.Logger

func Init() {
	level := logutil.Info
	switch config.GetConfig().LogLevel {
	case "debug":
		level = logutil.Debug
	case "info":
		level = logutil.Info
	case "warn":
		level = logutil.Warn
	case "error":
		level = logutil.Error
	}
	rp := pathutil.GetRootPath()
	logDir := filepath.Join(rp, "data")

	Logger = logutil.NewLogger(&logutil.Option{
		Dir:        logDir,
		Compress:   true,
		Level:      level,
		ReserveDay: 30,
	})
	return
}
