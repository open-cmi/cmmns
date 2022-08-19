package api

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
)

func ProcHello(c *gin.Context, msg json.RawMessage) {
	return
}

func TestRegisterAPI(t *testing.T) {
	RegisterAPIFunc("test", "Hello", ProcHello)

	proc := GetAPIFunc("test", "Hello")
	if proc == nil {
		t.Errorf("api get failed\n")
	}
}
