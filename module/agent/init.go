package agent

import (
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/service/ticker"
)

func init() {
	err := ticker.Register("agent-status-check", "0 */1 * * * *", func() {
		CheckStatus()
	})
	logger.Debug("agent register ticker, ", err)
}
