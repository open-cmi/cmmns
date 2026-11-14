package upgrademng

// 用于下载meta文件，通知用户有新版本，下载安装包后进行校验
// 这个文件是上传新版本时自动生成，用户不需要感知
type UpgradeMetaInfo struct {
	Prod    string `json:"prod"`
	Package string `json:"package"`
	SHA256  string `json:"sha256"` // sha256 hash
}
