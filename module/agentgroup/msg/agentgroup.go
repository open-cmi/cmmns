package msg

type CreateMsg struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MultiDeleteMsg struct {
	ID []string `json:"id"`
}

type EditMsg struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
