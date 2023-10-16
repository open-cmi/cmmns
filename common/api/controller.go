package api

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) map[string]interface{} {
	//sess, exist := c.Get("session")
	u, exist := c.Get("user")
	if exist {
		// session := sess.(*sessions.Session)
		// user, ok := session.Values["user"].(map[string]interface{})
		// if ok {
		// 	return user
		// }
		user, ok := u.(map[string]interface{})
		if ok {
			return user
		}
		return user
	}

	return nil
}

// ParseParams parse param
func ParseParams(c *gin.Context, option *Option) (err error) {
	var userID string = ""
	user := GetUser(c)
	if user != nil {
		userID = user["id"].(string)
	}
	option.UserID = userID
	option.Context = c

	devID := c.DefaultQuery("dev_id", "")
	option.DevID = devID

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
		option.Filters = []FilterOption{}
	}
	return err
}
