package prod

// NavConfig nav config
type NavConfig struct {
	Label string `json:"label"`
	Key   string `json:"key"`
	Link  string `json:"link"`
}

// NavTopConfig nav top config
type NavTopConfig struct {
	Label    string      `json:"label"`
	Key      string      `json:"key"`
	MenuList []NavConfig `json:"menulist"`
}

// BasicInfo prod info
type BasicInfo struct {
	Name string         `json:"name"`
	Nav  []NavTopConfig `json:"nav"`
}

var prodInfo BasicInfo

// GetProdInfo get prod nav from redis
func GetProdInfo() BasicInfo {
	return prodInfo
}

// SetProdInfo set prod info
func SetProdInfo(pi *BasicInfo) {
	prodInfo = *pi
	return
}
