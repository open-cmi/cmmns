package request

// 条件比较
type FilterQuery struct {
	Type      string      `json:"type"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Condition string      `json:"condition"`
}

// RequestQuery request param
type RequestQuery struct {
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Order    string        `json:"order"`
	OrderBy  string        `json:"orderby"`
	Filters  []FilterQuery `json:"filters"`
}
