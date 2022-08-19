package webserver

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// RouterGroup is a function to api group router
type RouterGroup func(e *gin.Engine)

var authGroup map[string]RouterGroup = make(map[string]RouterGroup)
var unauthGroup map[string]RouterGroup = make(map[string]RouterGroup)

// UnauthInit no auth router init
func UnauthInit(e *gin.Engine) {
	for _, apiGroup := range unauthGroup {
		apiGroup(e)
	}
}

// AuthInit auth router init
func AuthInit(e *gin.Engine) {

	for _, apiGroup := range authGroup {
		apiGroup(e)
	}
}

func RegisterAuthRouter(module string, group RouterGroup) error {
	_, found := authGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s auth group api has been registered", module)
		return errors.New(errMsg)
	}
	authGroup[module] = group
	return nil
}

func RegisterUnauthRouter(module string, group RouterGroup) error {
	_, found := unauthGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s unauth group api has been registered", module)
		return errors.New(errMsg)
	}
	unauthGroup[module] = group
	return nil
}
