package scheduler

type CreateMsg struct {
	JobType  string `json:"job_type"`
	Priority int    `json:"priority"`
	Content  string `json:"content"`
	RunType  string `json:"run_type"`
	RunSpec  string `json:"run_spec"`
}

type MultiDeleteMsg struct {
	ID []string `json:"id"`
}
