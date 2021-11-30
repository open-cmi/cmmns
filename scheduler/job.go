package scheduler

import "encoding/json"

// Request task request
type Request struct {
	JobID   string      `json:"jobid"`
	JobType string      `json:"jobtype"`
	Content interface{} `json:"content,omitempty"`
}

// Response task response
type Response struct {
	JobID   string          `json:"jobid"`
	JobType string          `json:"jobtype"`
	Results json.RawMessage `json:"results"`
}
