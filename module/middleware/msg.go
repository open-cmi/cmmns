package middleware

type CreateTokenRequest struct {
	Name      string `json:"name"`
	ExpireDay int    `json:"expire_day"`
}
