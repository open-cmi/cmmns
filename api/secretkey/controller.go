package secretkey

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/module/secretkey"
)

func NameList(c *gin.Context) {
	count, results, err := secretkey.NameList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": results,
		},
	})
}

func List(c *gin.Context) {
	var param goparam.Option
	goparam.ParseParams(c, &param)

	count, results, err := secretkey.List(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": results,
		},
	})
}

func MultiDelete(c *gin.Context) {
	var reqMsg secretkey.MultiDeleteMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := goparam.GetUser(c)
	userID, _ := user["id"].(string)

	var option goparam.Option
	option.UserID = userID
	err := secretkey.MultiDelete(&option, reqMsg.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func Create(c *gin.Context) {
	var reqMsg secretkey.CreateMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := goparam.GetUser(c)
	userID, _ := user["id"].(string)

	_, err := secretkey.Create(&goparam.Option{
		UserID: userID,
	}, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func Get(c *gin.Context) {
	identify := c.Param("id")

	user := goparam.GetUser(c)
	userID, _ := user["id"].(string)

	m := secretkey.Get(&goparam.Option{
		UserID: userID,
	}, identify)

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": m,
	})
}

func Delete(c *gin.Context) {
	identify := c.Param("id")

	user := goparam.GetUser(c)
	userID, _ := user["id"].(string)
	err := secretkey.Delete(&goparam.Option{
		UserID: userID,
	}, identify)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func Edit(c *gin.Context) {
	identify := c.Param("id")

	var reqMsg secretkey.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := goparam.GetUser(c)
	userID, _ := user["id"].(string)
	err := secretkey.Edit(&goparam.Option{
		UserID: userID,
	}, identify, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
