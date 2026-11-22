package syslog

import (
	"io"
	"os"
	"time"

	"github.com/nxadm/tail"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/ticker"
	"github.com/open-cmi/gobase/initial"
)

func watchFile(filePath string) {
	allowed := make(map[string]struct{}, len(gConf.Processes))
	for _, p := range gConf.Processes {
		if p != "" {
			allowed[p] = struct{}{}
		}
	}

	for {
		t, err := tail.TailFile(filePath, tail.Config{
			Follow:    true,
			ReOpen:    true,
			MustExist: true,
			Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
			Logger:    tail.DiscardingLogger,
		})
		if err != nil {
			logger.Errorf("tail syslog file %s failed: %s\n", filePath, err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		for line := range t.Lines {
			if line.Err != nil {
				logger.Debugf("tail syslog line error: %s\n", line.Err.Error())
				continue
			}
			m, ok := ParseLine(line.Text)
			if !ok {
				continue
			}
			if len(allowed) > 0 {
				if _, exists := allowed[m.Process]; !exists {
					continue
				}
			}
			if err := m.Save(); err != nil {
				logger.Errorf("save syslog record failed: %s\n", err.Error())
			}
		}

		_ = t.Wait()
		t.Cleanup()
		time.Sleep(5 * time.Second)
	}
}

var defaultSyslogFiles = []string{
	"/var/log/syslog",
	"/var/log/messages",
}

func resolveSyslogFile() string {
	if gConf.FilePath != "" {
		if _, err := os.Stat(gConf.FilePath); err == nil {
			return gConf.FilePath
		}
		return ""
	}

	for _, path := range defaultSyslogFiles {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func watchSyslog() {
	if watchJournal() {
		return
	}

	if filePath := resolveSyslogFile(); filePath != "" {
		watchFile(filePath)
		return
	}

	logger.Debugf("neither systemd journal nor syslog file found, skip syslog watcher\n")
}

// Init starts the background syslog watcher on service startup
func Init() error {
	go watchSyslog()
	return nil
}

func init() {
	initial.Register("syslog", initial.PhaseDefault, Init)

	// Clean expired syslog daily at 02:00:00
	ticker.Register("syslog-clean-ticker", "0 0 2 * * *", func(name string, data interface{}) {
		keepDays := gConf.KeepDays
		if keepDays <= 0 {
			keepDays = 30
		}
		_ = CleanExpired(keepDays)
	}, nil)
}
