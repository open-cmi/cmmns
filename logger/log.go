package logger

import (
	"path/filepath"

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
	rp := pathutil.GetRootPath()
	logDir := filepath.Join(rp, "data")

	Logger = logutil.NewLogger(&logutil.Option{
		Dir:        logDir,
		Compress:   true,
		Level:      logutil.Info,
		ReserveDay: 30,
	})
	return
}
