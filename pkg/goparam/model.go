package goparam

import (
	"fmt"

	"github.com/open-cmi/cmmns/pkg/structure"
)

// 条件比较
type FilterParam struct {
	Type      string      `json:"type"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Condition string      `json:"condition"`
}

// PageParam page param
type PageParam struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// OrderParam order option
type OrderParam struct {
	Order   string `json:"order"`
	OrderBy string `json:"order_by"`
}

// Param model option
type Param struct {
	UserID     string
	Role       string
	DevID      string
	PageParam  PageParam
	OrderParam OrderParam
	Filters    []FilterParam
}

func GetColumn(v interface{}, skipColumn []string) []string {

	var skipmap map[string]bool = make(map[string]bool, 0)
	for _, sc := range skipColumn {
		skipmap[sc] = true
	}

	var fields []string
	dbFields := structure.GetStructFields(v, "db")
	for _, field := range dbFields {
		if !skipmap[field] {
			fields = append(fields, field)
		}
	}
	return fields
}

func GetColumnUpdateNamed(columns []string) []string {
	var updates []string = []string{}
	for _, column := range columns {
		updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
	}
	return updates
}

func GetColumnInsertNamed(columns []string) []string {
	var named []string = []string{}
	for _, col := range columns {
		seq := fmt.Sprintf(`:%s`, col)
		named = append(named, seq)
	}
	return named
}

// BuildWhereClause build where clause
func BuildWhereClause(opt *Param) (format string, args []interface{}) {
	if opt == nil {
		return "", []interface{}{}
	}
	var clause string = ""

	args = []interface{}{}
	for index, filter := range opt.Filters {
		if index == 0 {
			clause += " where"
		} else {
			clause += " and"
		}

		if filter.Type == "string" {
			value := filter.Value.(string)
			if filter.Condition == "contains" {
				clause += fmt.Sprintf(` %s like '%%' || $%d || '%%'`, filter.Name, index+1)
				args = append(args, value)
			} else if filter.Condition == "eq" {
				clause += fmt.Sprintf(" %s = $%d", filter.Name, index+1)
				args = append(args, value)
			}
		} else if filter.Type == "number" || filter.Type == "integer" {
			value := filter.Value.(int)
			if filter.Condition == "eq" {
				clause += fmt.Sprintf(" %s = $%d", filter.Name, index+1)
				args = append(args, value)
			} else if filter.Condition == "lt" {
				clause += fmt.Sprintf(" %s < $%d", filter.Name, index+1)
				args = append(args, value)
			} else if filter.Condition == "gt" {
				clause += fmt.Sprintf(" %s > $%d", filter.Name, index+1)
				args = append(args, value)
			} else if filter.Condition == "le" {
				clause += fmt.Sprintf(" %s <= $%d", filter.Name, index+1)
				args = append(args, value)
			} else if filter.Condition == "ge" {
				clause += fmt.Sprintf(" %s >= $%d", filter.Name, index+1)
				args = append(args, value)
			}
		}
	}

	return clause, args
}

// BuildFinalClause build final clause
func BuildFinalClause(opt *Param) string {
	if opt == nil {
		return ""
	}

	var clause string = ""

	if opt.OrderParam.OrderBy != "" {
		if opt.OrderParam.Order == "" {
			opt.OrderParam.Order = "asc"
		}
		clause += fmt.Sprintf(` ORDER BY %s %s`, opt.OrderParam.OrderBy, opt.OrderParam.Order)
	}

	// limit must before offset
	if opt.PageParam.PageSize != 0 {
		offset := opt.PageParam.Page * opt.PageParam.PageSize
		clause += fmt.Sprintf(" LIMIT %d OFFSET %d", opt.PageParam.PageSize, offset)
	}

	return clause
}
