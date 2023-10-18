package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/email"
)

func GetEmail(c *gin.Context) {
	m := email.Get()
	if m == nil {
		m = email.New()
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "", "data": *m})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": *m})
}

func SetEmail(c *gin.Context) {
	var req email.SetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := email.SetSenderEmail(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}
