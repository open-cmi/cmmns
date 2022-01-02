package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	msg "github.com/open-cmi/cmmns/msg/request"
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
	} else {
		p.Filters = []msg.FilterQuery{}
	}
	return err
}

func BuildWhereClause(r *msg.RequestQuery) (format string, args []interface{}) {
	var clause string = ""

	args = []interface{}{}
	if len(r.Filters) != 0 {
		for index, filter := range r.Filters {
			if index == 0 {
				clause += " where"
			} else {
				clause += " and"
			}

			if filter.Type == "string" {
				value := filter.Value.(string)
				if filter.Condition == "contains" {
					clause += fmt.Sprintf(" %s like $%d", filter.Name, index+1)
					args = append(args, value)
				} else if filter.Condition == "eq" {
					clause += fmt.Sprintf(" %s = $%d", filter.Name, index+1)
					args = append(args, value)
				}
			} else if filter.Type == "number" {
				value := filter.Value.(int32)
				if filter.Condition == "eq" {
					clause += fmt.Sprintf(" %s = $%d", filter.Name, index+1)
					args = append(args, value)
				} else if filter.Condition == "lt" {
					clause += fmt.Sprintf(" %s < $%d", filter.Name, index+1)
					args = append(args, value)
				} else if filter.Condition == "gt" {
					clause += fmt.Sprintf(" %s > $%d", filter.Name, index+1)
					args = append(args, value)
				}
			}
		}
	}

	return clause, args
}

func BuildFinalClause(r *msg.RequestQuery) string {
	var clause string = ""

	if r.OrderBy != "" && r.Order != "" {
		clause += fmt.Sprintf(` ORDER BY %s %s`, r.OrderBy, r.Order)
	}
	offset := r.Page * r.PageSize
	clause += fmt.Sprintf(" OFFSET %d LIMIT %d", offset, r.PageSize)

	return clause
}
