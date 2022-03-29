package setting

type EditMsg struct {
	Scope  string `json:"scope"`
	Belong string `json:"belong"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
