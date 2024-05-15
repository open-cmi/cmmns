package time

import (
	"github.com/open-cmi/cmmns/service/business"
)

func GetTimeZoneList() ([]string, error) {
	return []string{"Asia/Shanghai"}, nil
}

func SetTimeZone(tz string) error {
	return nil
}

func Init() error {
	m := Get()
	if m == nil {
		m = New()
		m.AutoAdjust = true
		m.NtpServer = "ntp.ubuntu.com"
		m.TimeZone = "Asia/Shanghai"
		return m.Save()
	}
	return nil
}

func init() {
	business.Register("time-setting", business.DefaultPriority, Init)
}
