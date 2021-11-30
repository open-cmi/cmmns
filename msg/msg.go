package msg

// RequestParams request param
type RequestParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Order    string `json:"order"`
	OrderBy  string `json:"orderby"`
}
