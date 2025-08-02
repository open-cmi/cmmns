package captcha

import (
	"bytes"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// GetID get captcha id
func GetID(c *gin.Context) {
	captchaid := captcha.NewLen(4)

	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"captchaid": captchaid,
		},
	})
}

// GetPic get pick accord by id
func GetPic(c *gin.Context) {
	var content bytes.Buffer
	c.Header("Content-Type", "image/png")
	id := c.Param("id")
	captcha.WriteImage(&content, id, 120, 56)
	reader := bytes.NewReader(content.Bytes())

	c.DataFromReader(http.StatusOK, int64(len(content.Bytes())), "image/png", reader, nil)
}
