package licmng

type LicenseInfo struct {
	ExpireTime int64    `json:"expire_time"`
	Version    string   `json:"version"`
	Modules    []string `json:"modules"`
	MCode      string   `json:"mcode,omitempty"`
}

type CreateLicenseRequest struct {
	Customer   string `json:"customer"`
	Prod       string `json:"prod"`
	ExpireTime int64  `json:"expire_time"`
	MCode      string `json:"mcode,omitempty"`
	Modules    string `json:"modules"`
	Version    string `json:"version"`
}
