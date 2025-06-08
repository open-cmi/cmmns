package time

import (
	"github.com/open-cmi/cmmns/initial"
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
		m.NtpServer = "cn.ntp.org.cn"
		m.TimeZone = "Asia/Shanghai"
		return m.Save()
	}
	return nil
}

func init() {
	initial.Register("time-setting", initial.DefaultPriority, Init)
}
