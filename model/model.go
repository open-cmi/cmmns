package model

import (
	"fmt"
	"reflect"
)

// 条件比较
type FilterOption struct {
	Type      string      `json:"type"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Condition string      `json:"condition"`
}

// PageOption page param
type PageOption struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// OrderOption order option
type OrderOption struct {
	Order   string `json:"order"`
	OrderBy string `json:"orderby"`
}

// Option model option
type Option struct {
	UserID      string
	PageOption  PageOption
	OrderOption OrderOption
	Filters     []FilterOption
}

type Model struct {
	ID          string `json:"id" db:"id"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
	UpdatedTime int64  `json:"updated_time" db:"updated_time"`
}

func GetColumn(v interface{}, skipColumn []string) []string {
	var columns []string = []string{}
	dataValue := reflect.ValueOf(v)
	if dataValue.Kind() != reflect.Struct {
		return columns
	}

	var skipmap map[string]bool = make(map[string]bool, 0)
	for _, sc := range skipColumn {
		skipmap[sc] = true
	}

	t := dataValue.Type()
	for i := 0; i < t.NumField(); i++ {
		field := dataValue.Type().Field(i)
		tagColumn, ok := field.Tag.Lookup("db")
		if ok {
			sk, _ := skipmap[tagColumn]
			if !sk {
				columns = append(columns, tagColumn)
			}
		}
	}

	return columns
}

func GetColumnNamed(columns []string) []string {
	var named []string = []string{}
	for _, col := range columns {
		seq := fmt.Sprintf(`:%s`, col)
		named = append(named, seq)
	}
	return named
}

// BuildWhereClause build where clause
func BuildWhereClause(o *Option) (format string, args []interface{}) {
	var clause string = ""

	var whereUsed bool = false

	args = []interface{}{}
	for index, filter := range o.Filters {
		if whereUsed {
			clause += " and"
		} else {
			clause += " where"
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

	return clause, args
}

// BuildFinalClause build final clause
func BuildFinalClause(o *Option) string {
	var clause string = ""

	if o.OrderOption.OrderBy != "" {
		if o.OrderOption.Order == "" {
			o.OrderOption.Order = "asc"
		}
		clause += fmt.Sprintf(` ORDER BY %s %s`, o.OrderOption.OrderBy, o.OrderOption.Order)
	}

	offset := o.PageOption.Page * o.PageOption.PageSize
	clause += fmt.Sprintf(" LIMIT %d OFFSET %d", o.PageOption.PageSize, offset)

	return clause
}
