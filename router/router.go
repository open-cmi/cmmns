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
}

// AuthInit auth router init
func AuthInit(e *gin.Engine) {
	fmt.Println("AuthInit")
	UserAuthGroup(e)
	SystemInfoGroup(e)
}
