package middleware

type CreateTokenRequest struct {
	Name      string `json:"name"`
	ExpireDay int    `json:"expire_day"`
}

type DeleteTokenRequest struct {
	Name string `json:"name"`
}
