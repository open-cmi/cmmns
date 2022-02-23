package manhour

type CreateMsg struct {
	Date      int64  `json:"date"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Content   string `json:"content"`
}

type MultiDeleteMsg struct {
	ID []string `json:"id"`
}

type EditMsg struct {
	Date      int64  `json:"date"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Content   string `json:"content"`
}
