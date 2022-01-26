package router

import (
	"github.com/gin-gonic/gin"
)

// NauthInit no auth router init
func NauthInit(e *gin.Engine) {
	CaptchaGroup(e)
	UserNauthGroup(e)
	AgentNauthGroup(e)
}

// AuthInit auth router init
func AuthInit(e *gin.Engine) {
	AgentAuthGroup(e)
	MasterAuthGroup(e)
	AgentGroupAuthGroup(e)
	UserAuthGroup(e)
	SystemGroup(e)
	ProgressAuthGroup(e)
	AuditLogAuthGroup(e)
	AssistAuthGroup(e)
	SecretKeyAuthGroup(e)
	TemplateAuthGroup(e)
}
