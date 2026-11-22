//go:build linux

package syslog

import (
	"strconv"
	"time"

	"github.com/coreos/go-systemd/v22/sdjournal"
	"github.com/google/uuid"
	"github.com/open-cmi/gobase/essential/logger"
)

func watchJournal() bool {
	j, err := sdjournal.NewJournal()
	if err != nil {
		logger.Debugf("open systemd journal failed: %s\n", err.Error())
		return false
	}
	defer j.Close()

	for _, proc := range gConf.Processes {
		if proc != "" {
			if err := j.AddMatch("SYSLOG_IDENTIFIER=" + proc); err != nil {
				logger.Errorf("add journal match %s failed: %s\n", proc, err.Error())
			}
		}
	}

	if err := j.SeekTail(); err != nil {
		logger.Errorf("seek journal tail failed: %s\n", err.Error())
		return false
	}
	_, _ = j.Previous()

	for {
		n, err := j.Next()
		if err != nil {
			logger.Errorf("journal next failed: %s\n", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		if n == 0 {
			j.Wait(sdjournal.IndefiniteWait)
			continue
		}

		entry, err := j.GetEntry()
		if err != nil {
			continue
		}

		proc := entry.Fields["SYSLOG_IDENTIFIER"]
		if proc == "" {
			proc = entry.Fields["_COMM"]
		}
		msg := entry.Fields["MESSAGE"]
		if msg == "" {
			continue
		}

		pid, _ := strconv.Atoi(entry.Fields["_PID"])
		level := priorityLevelMap[entry.Fields["PRIORITY"]]
		now := time.Now().Unix()
		logTime := int64(entry.RealtimeTimestamp / 1_000_000)
		if logTime == 0 {
			logTime = now
		}

		m := &Model{
			ID:        uuid.New().String(),
			LogTime:   logTime,
			Process:   proc,
			PID:       pid,
			Hostname:  entry.Fields["_HOSTNAME"],
			Level:     level,
			Message:   msg,
			CreatedAt: now,
		}
		if err := m.Save(); err != nil {
			logger.Errorf("save journal syslog record failed: %s\n", err.Error())
		}
	}
}
