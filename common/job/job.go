package job

import "encoding/json"

// Proc job proc callback
type Proc func(req *Request, resp *Response)

type Request struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content,omitempty"`
}

type Response struct {
	ID     string      `json:"id"`
	Type   string      `json:"type"`
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result,omitempty"`
}
