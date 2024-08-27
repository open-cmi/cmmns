package upgrademng

type UpgradeMetaInfo struct {
	Prod    string `json:"prod"`
	Package string `json:"package"`
	SHA256  string `json:"sha256"` // sha256 hash
}
