package def

// JobRequest 用户添加，并且需要返回给agent用于执行
type JobRequest struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
}

// JobResponse agent返回，并且返回给用户
type JobResponse struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result string `json:"result,omitempty"`
}
