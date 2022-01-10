package log

import (
	"path/filepath"

	"github.com/open-cmi/goutils/common"
	"github.com/open-cmi/goutils/logutils"
)

const (
	Debug = iota
	Info
	Warn
	Error
)

var Logger *logutils.Logger

func Init() {
	rp := common.GetRootPath()
	logDir := filepath.Join(rp, "data")

	Logger = logutils.NewLogger(&logutils.Option{
		Dir:        logDir,
		Compress:   true,
		Level:      logutils.Info,
		ReserveDay: 30,
	})
	return
}
