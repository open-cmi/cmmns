package auth

// User auth user struct
type User struct {
	ID  string `json:"id"`
	PID string `json:"parentid,omitempty"`
}
