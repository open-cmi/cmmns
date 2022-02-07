package scheduler

// Job job
type Job struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Priority int    `json:"priority"`
	State    string `json:"-"`
	Count    int    `json:"-"`
	Content  string `json:"content,omitempty"`
}
