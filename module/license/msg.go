package license

type LicenseInfo struct {
	ExpireTime int64    `json:"expire_time"`
	Version    string   `json:"version"`
	Modules    []string `json:"modules"`
}
