package logger

import (
	"fmt"

	"path/filepath"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/goutils/logutil"
	"github.com/open-cmi/goutils/pathutil"
)

type Feature interface {
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

var Logger Feature

type Config struct {
	Level string `json:"level"`
	Path  string `json:"path,omitempty"`
}

var moduleConfig Config

func RegisterLogger(l interface{}) {
	Logger = l.(Feature)
}

func Init() error {
	level := logutil.Info
	switch moduleConfig.Level {
	case "debug":
		level = logutil.Debug
	case "info":
		level = logutil.Info
	case "warn":
		level = logutil.Warn
	case "error":
		level = logutil.Error
	}
	logPath := moduleConfig.Path
	if logPath == "" {
		rp := pathutil.GetRootPath()
		logPath = filepath.Join(rp, "data")
	}

	Logger = logutil.NewLogger(&logutil.Option{
		Dir:        logPath,
		Compress:   true,
		Level:      level,
		ReserveDay: 30,
	})

	return nil
}

func Error(v ...interface{}) {
	if Logger != nil {
		Logger.Error(v...)
	} else {
		fmt.Println(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Errorf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

func Warn(v ...interface{}) {
	if Logger != nil {
		Logger.Warn(v...)
	} else {
		fmt.Println(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Warnf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

func Info(v ...interface{}) {
	if Logger != nil {
		Logger.Info(v...)
	} else {
		fmt.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Infof(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if Logger != nil {
		Logger.Debug(v...)
	} else {
		fmt.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Debugf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

func init() {
	moduleConfig.Level = "debug"
	moduleConfig.Path = "/tmp/"
	config.RegisterConfig("log", &moduleConfig)
}
