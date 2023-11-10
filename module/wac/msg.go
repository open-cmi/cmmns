package wac

type SetRequest struct {
	Enable       bool   `json:"enable"`
	Mode         string `json:"mode"`
	RawWhitelist string `json:"raw_whitelist"`
	RawBlacklist string `json:"raw_blacklist"`
}
