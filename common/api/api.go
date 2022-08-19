package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type MsgProc func(c *gin.Context, msg json.RawMessage)

var APIMapping map[string]map[string]MsgProc = make(map[string]map[string]MsgProc)

func GetAPIFunc(prod string, msgType string) MsgProc {
	prodMapping, found := APIMapping[prod]
	if !found {
		return nil
	}
	return prodMapping[msgType]
}

func RegisterAPIFunc(prod string, msgType string, proc MsgProc) error {
	prodMapping, found := APIMapping[prod]
	if !found {
		APIMapping[prod] = make(map[string]MsgProc)
		prodMapping = APIMapping[prod]
	}
	_, found = prodMapping[msgType]
	if found {
		errMsg := fmt.Sprintf("%s msg type %s has been registered", prod, msgType)
		return errors.New(errMsg)
	}
	prodMapping[msgType] = proc
	return nil
}
