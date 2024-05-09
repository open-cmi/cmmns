package licmng

type LicenseInfo struct {
	ExpireTime int64    `json:"expire_time"`
	Version    string   `json:"version"`
	Modules    []string `json:"modules"`
	Machine    string   `json:"machine,omitempty"`
}

type CreateLicenseRequest struct {
	Customer   string `json:"customer"`
	Prod       string `json:"prod"`
	ExpireTime int64  `json:"expire_time"`
	Machine    string `json:"machine,omitempty"`
	Modules    string `json:"modules"`
	Version    string `json:"version"`
}
