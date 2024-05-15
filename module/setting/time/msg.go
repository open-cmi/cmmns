package time

type SettingRequest struct {
	TimeZone   string `json:"timezone"`
	NtpServer  string `json:"ntp_server"`
	AutoAdjust bool   `json:"auto_adjust"`
}
