package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	msg "github.com/open-cmi/cmmns/msg/common"
)

// ParseParams parse param
func ParseParams(c *gin.Context, p *msg.RequestQuery) (err error) {
	pagestr := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		page = 0
	}
	p.Page = page

	// page size
	pagesizestr := c.DefaultQuery("page_size", "25")
	pagesize, err := strconv.Atoi(pagesizestr)
	if err != nil {
		pagesize = 25
	}
	p.PageSize = pagesize

	p.Order = c.Query("order")
	p.OrderBy = c.Query("orderBy")

	filters := c.Query("filters")
	if filters != "" {
		err = json.Unmarshal([]byte(filters), &p.Filters)
		// 记录日志
	}
	return err
}

func BuildSQLClause(r *msg.RequestQuery) string {
	var clause string = ""
	if len(r.Filters) != 0 {
		for index, filter := range r.Filters {
			if index == 0 {
				clause += " where"
			} else {
				clause += " and"
			}

			if filter.Condition == "like" {
				clause += fmt.Sprintf(" %s like '%%%s%%'", filter.Key, filter.Value)
			}
		}
	}
	return clause
}
