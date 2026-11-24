package sshlogin

import (
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nxadm/tail"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/initial"
)

var defaultLogs = []string{"/var/log/auth.log", "/var/log/secure"}
var gLogFile = ""

var (
	gDedupMutex sync.Mutex
	gDedupCache = make(map[string]int64)
)

func extractUserHostPort(rest string) (string, string, int) {
	fromIdx := strings.Index(rest, " from ")
	if fromIdx == -1 {
		return "", "", 0
	}

	user := strings.TrimSpace(rest[:fromIdx])
	afterFrom := strings.TrimSpace(rest[fromIdx+len(" from "):])

	portIdx := strings.Index(afterFrom, " port ")
	if portIdx == -1 {
		rhost, _, _ := strings.Cut(afterFrom, " ")
		return user, strings.TrimSpace(rhost), 0
	}

	rhost := strings.TrimSpace(afterFrom[:portIdx])
	afterPort := strings.TrimSpace(afterFrom[portIdx+len(" port "):])
	portStr, _, _ := strings.Cut(afterPort, " ")
	port, _ := strconv.Atoi(portStr)

	return user, rhost, port
}

func isDuplicateEvent(key string, now int64) bool {
	gDedupMutex.Lock()
	defer gDedupMutex.Unlock()

	for k, ts := range gDedupCache {
		if now-ts > 5 {
			delete(gDedupCache, k)
		}
	}

	if lastTs, exists := gDedupCache[key]; exists && now-lastTs <= 2 {
		return true
	}
	gDedupCache[key] = now
	return false
}

// ParseSSHLogLine 解析 sshd 输出的 syslog 日志行
func ParseSSHLogLine(origin string) *SSHLoginEvent {
	trimmed := strings.TrimSpace(origin)
	if !strings.Contains(trimmed, "sshd") {
		return nil
	}

	now := time.Now().Unix()

	// 1. 登录成功匹配
	if idx := strings.Index(trimmed, "Accepted "); idx != -1 {
		rest := trimmed[idx+len("Accepted "):]
		authType := "other"
		if strings.HasPrefix(rest, "password for ") {
			authType = "password"
			rest = strings.TrimPrefix(rest, "password for ")
		} else if strings.HasPrefix(rest, "publickey for ") {
			authType = "publickey"
			rest = strings.TrimPrefix(rest, "publickey for ")
		} else if strings.HasPrefix(rest, "keyboard-interactive/pam for ") {
			authType = "keyboard-interactive"
			rest = strings.TrimPrefix(rest, "keyboard-interactive/pam for ")
		} else if strings.HasPrefix(rest, "keyboard-interactive for ") {
			authType = "keyboard-interactive"
			rest = strings.TrimPrefix(rest, "keyboard-interactive for ")
		} else if strings.Contains(rest, " for ") {
			forIdx := strings.Index(rest, " for ")
			authType = strings.TrimSpace(rest[:forIdx])
			rest = rest[forIdx+len(" for "):]
		}

		user, rhost, port := extractUserHostPort(rest)
		if user == "" || rhost == "" {
			return nil
		}

		dedupKey := "success:" + rhost + ":" + user
		if isDuplicateEvent(dedupKey, now) {
			return nil
		}

		return &SSHLoginEvent{
			Timestamp: now,
			Username:  user,
			RHost:     rhost,
			Port:      port,
			Status:    "success",
			AuthType:  authType,
		}
	}

	// 2. 登录失败匹配
	if idx := strings.Index(trimmed, "Failed "); idx != -1 {
		rest := trimmed[idx+len("Failed "):]
		authType := "other"
		if strings.HasPrefix(rest, "password for ") {
			authType = "password"
			rest = strings.TrimPrefix(rest, "password for ")
		} else if strings.HasPrefix(rest, "publickey for ") {
			authType = "publickey"
			rest = strings.TrimPrefix(rest, "publickey for ")
		} else if strings.Contains(rest, " for ") {
			forIdx := strings.Index(rest, " for ")
			authType = strings.TrimSpace(rest[:forIdx])
			rest = rest[forIdx+len(" for "):]
		}

		isInvalid := 0
		failReason := "bad_password"
		if authType == "publickey" {
			failReason = "bad_publickey"
		}
		if strings.HasPrefix(rest, "invalid user ") {
			isInvalid = 1
			failReason = "invalid_user"
			rest = strings.TrimPrefix(rest, "invalid user ")
		}

		user, rhost, port := extractUserHostPort(rest)
		if user == "" || rhost == "" {
			return nil
		}

		dedupKey := "failed:" + rhost + ":" + user
		if isDuplicateEvent(dedupKey, now) {
			return nil
		}

		return &SSHLoginEvent{
			Timestamp:  now,
			Username:   user,
			RHost:      rhost,
			Port:       port,
			Status:     "failed",
			AuthType:   authType,
			IsInvalid:  isInvalid,
			FailReason: failReason,
		}
	}

	// 3. 兜底匹配 PAM 认证失败
	if strings.Contains(trimmed, "authentication failure;") && strings.Contains(trimmed, "rhost=") {
		rhostIndex := strings.Index(trimmed, "rhost=")
		text := strings.TrimPrefix(trimmed[rhostIndex:], "rhost=")
		rhost, _, _ := strings.Cut(text, " ")
		rhost = strings.TrimSpace(rhost)
		if rhost == "" {
			return nil
		}

		user := "unknown"
		if userIndex := strings.Index(trimmed, " user="); userIndex != -1 {
			uText := trimmed[userIndex+len(" user="):]
			user, _, _ = strings.Cut(uText, " ")
			user = strings.TrimSpace(user)
		}

		dedupKey := "failed:" + rhost + ":" + user
		if isDuplicateEvent(dedupKey, now) {
			return nil
		}

		return &SSHLoginEvent{
			Timestamp:  now,
			Username:   user,
			RHost:      rhost,
			Port:       0,
			Status:     "failed",
			AuthType:   "pam",
			FailReason: "pam_failure",
		}
	}

	return nil
}

func tailFile(filePath string) {
	for {
		t, err := tail.TailFile(filePath, tail.Config{
			Follow:    true,
			ReOpen:    true,
			MustExist: false,
			Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
			Logger:    tail.DiscardingLogger,
		})
		if err != nil {
			logger.Errorf("Error tailing file %s: %s\n", filePath, err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		for line := range t.Lines {
			if line == nil {
				continue
			}
			if line.Err != nil {
				logger.Warnf("Error reading line from %s: %v", filePath, line.Err)
				continue
			}
			event := ParseSSHLogLine(line.Text)
			if event != nil {
				HandleSSHLoginEvent(event)
			}
		}

		_ = t.Stop()
		t.Cleanup()
		time.Sleep(1 * time.Second)
	}
}

// Init 启动 SSH 登录日志监控服务
func Init() error {
	if gLogFile == "" {
		logger.Warnf("ssh login monitor: no auth log file found in %v", defaultLogs)
		return nil
	}
	logger.Infof("ssh login monitor started on %s", gLogFile)
	go tailFile(gLogFile)
	return nil
}

func init() {
	for _, lf := range defaultLogs {
		_, err := os.Stat(lf)
		if err != nil {
			continue
		}
		gLogFile = lf
		break
	}

	initial.Register("ssh-login-monitor", initial.PhaseDefault, Init)
}
