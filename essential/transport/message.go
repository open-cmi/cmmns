package transport

import "encoding/json"

// Response msg
type Response struct {
	Ret  int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data,omitempty"`
}
