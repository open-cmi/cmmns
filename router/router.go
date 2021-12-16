package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// NauthInit no auth router init
func NauthInit(e *gin.Engine) {
	fmt.Println("NauthInit")

	CaptchaGroup(e)
	UserNauthGroup(e)
	AgentNauthGroup(e)
}

// AuthInit auth router init
func AuthInit(e *gin.Engine) {
	fmt.Println("AuthInit")
	AgentAuthGroup(e)
	AgentGroupGroup(e)
	UserAuthGroup(e)
	SystemInfoGroup(e)
	ProgressAuthGroup(e)
	AuditLogAuthGroup(e)
}
