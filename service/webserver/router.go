package webserver

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// APIGroup is a function to api group router
type APIGroup func(e *gin.Engine)

var authGroup map[string]APIGroup = make(map[string]APIGroup)
var unauthGroup map[string]APIGroup = make(map[string]APIGroup)

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

func RegisterAuthAPI(module string, group APIGroup) error {
	_, found := authGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s auth group api has been registered", module)
		return errors.New(errMsg)
	}
	authGroup[module] = group
	return nil
}

func RegisterUnauthAPI(module string, group APIGroup) error {
	_, found := unauthGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s unauth group api has been registered", module)
		return errors.New(errMsg)
	}
	unauthGroup[module] = group
	return nil
}
