package agent

import (
	"time"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/service/ticker"
)

func CheckStatus() {
	if !gConf.CheckStatus {
		return
	}

	logger.Debugf("start to check agent status\n")

	_, all, err := List(nil)
	if err != nil {
		logger.Errorf("list agent failed: %s\n", err.Error())
		return
	}

	now := time.Now().Unix()
	for _, item := range all {
		if now-item.UpdatedTime > int64(5*time.Minute/time.Second) {
			item.State = StateOffline
			item.Save()
		}
	}
}

func init() {
	err := ticker.Register("agent-status-check", "0 */5 * * * *", func() {
		CheckStatus()
	})
	if err != nil {
		logger.Errorf("register ticker failed: %s\n", err.Error())
	}
}
