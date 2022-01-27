package controller

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/model"
)

func GetUser(c *gin.Context) map[string]interface{} {
	sess, _ := c.Get("session")
	session := sess.(*sessions.Session)
	user, ok := session.Values["user"].(map[string]interface{})
	if ok {
		return user
	}
	return nil
}

// ParseParams parse param
func ParseParams(c *gin.Context, option *model.Option) (err error) {
	var userID string = ""
	user := GetUser(c)
	if user != nil {
		userID = user["id"].(string)
	}
	option.UserID = userID

	pagestr := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		page = 0
	}
	option.PageOption.Page = page

	// page size
	pagesizestr := c.DefaultQuery("page_size", "25")
	pagesize, err := strconv.Atoi(pagesizestr)
	if err != nil {
		pagesize = 25
	}
	option.PageOption.PageSize = pagesize

	option.OrderOption.Order = c.Query("order")
	option.OrderOption.OrderBy = c.Query("orderBy")

	filters := c.Query("filters")
	if filters != "" {
		err = json.Unmarshal([]byte(filters), &option.Filters)
		// 记录日志
	} else {
		option.Filters = []model.FilterOption{}
	}
	return err
}
