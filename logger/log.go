package logger

import (
	"path/filepath"

	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/goutils/logutil"
	"github.com/open-cmi/goutils/pathutil"
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

/*
func Error(v ...interface{}) {
	Logger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	Logger.Errorf(format, v...)
}

func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	Logger.Warnf(format, v...)
}

func Info(v ...interface{}) {
	Logger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	Logger.Infof(format, v...)
}

func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	Logger.Debugf(format, v...)
}
*/
