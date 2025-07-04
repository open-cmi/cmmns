package time

import "github.com/open-cmi/gobase/initial"

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
	initial.Register("time", initial.PhaseDefault, Init)
}
