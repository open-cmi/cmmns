package licmng

type LicenseInfo struct {
	ExpireTime int64    `json:"expire_time"`
	Version    string   `json:"version"`
	Modules    []string `json:"modules"`
	MCode      string   `json:"mcode,omitempty"`
}

type CreateLicenseRequest struct {
	Customer    string `json:"customer"`
	Prod        string `json:"prod"`
	ValidPeriod int    `json:"valid_period"`
	MCode       string `json:"mcode,omitempty"`
	Modules     string `json:"modules"`
	Version     string `json:"version"`
	Model       string `json:"model"`
}
