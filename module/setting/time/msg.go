package time

type SettingRequest struct {
	NtpServer  string `json:"ntp_server"`
	AutoAdjust bool   `json:"auto_adjust"`
}
