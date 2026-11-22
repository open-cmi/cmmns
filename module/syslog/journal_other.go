//go:build !linux

package syslog

import (
	"github.com/open-cmi/gobase/essential/logger"
)

func watchJournal() bool {
	logger.Debugf("systemd journal watcher is only supported on linux, skip\n")
	return false
}
