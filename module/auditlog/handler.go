package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

type AuditLogHandler struct {
	Username string
	IP       string
}

func (h *AuditLogHandler) InsertLoginLog(username string, action string, success bool) error {
	return InsertLog(h.IP, username, LoginType, action, success)
}

func (h *AuditLogHandler) InsertOperationLog(action string, success bool) error {
	return InsertLog(h.IP, h.Username, OperationType, action, success)
}

func NewAuditHandler(c *gin.Context) *AuditLogHandler {
	ip := c.ClientIP()

	var username string
	usermap := goparam.GetUser(c)
	if usermap != nil {
		username, _ = usermap["username"].(string)
	}
	return &AuditLogHandler{
		IP:       ip,
		Username: username,
	}
}
