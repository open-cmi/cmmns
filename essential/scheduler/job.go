package scheduler

// Job job
type Job struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Priority int         `json:"priority"`
	State    string      `json:"state"`
	Count    int         `json:"count"`
	Done     int         `json:"done"`
	Content  string      `json:"content,omitempty"`
	Result   interface{} `json:"result,omitempty"`
}
