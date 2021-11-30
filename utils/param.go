package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/msg"
)

// ParseParams parse param
func ParseParams(c *gin.Context, p *msg.RequestParams) {
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
		pagesize = 0
	}
	p.PageSize = pagesize

	p.Order = c.Query("order")
	p.OrderBy = c.Query("orderBy")
}
