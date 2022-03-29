package template

type CreateMsg struct {
	Name string `json:"name"`
}

type MultiDeleteMsg struct {
	Name []string `json:"name"`
}

type EditMsg struct {
	Name string `json:"name"`
}
