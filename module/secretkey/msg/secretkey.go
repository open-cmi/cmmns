package msg

type CreateMsg struct {
	Name         string `json:"name"`
	KeyType      string `json:"key_type"`
	KeyLength    int    `json:"key_length"`
	Comment      string `json:"comment"`
	PassPhrase   string `json:"passphrase"`
	Confirmation string `json:"confirmation"`
}

type MultiDeleteMsg struct {
	Name []string `json:"name"`
}

type EditMsg struct {
	Name         string `json:"name"`
	KeyType      string `json:"key_type"`
	KeyLength    int    `json:"key_length"`
	Comment      string `json:"comment"`
	PassPhrase   string `json:"passphrase"`
	Confirmation string `json:"confirmation"`
}
