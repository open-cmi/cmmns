package wac

import "github.com/open-cmi/cmmns/service/business"

func Init() error {
	m := Get()
	if m == nil {
		m = New()
		m.Mode = "blacklist"
	}
	gWebAccessControl = m.ConvertoWAC()
	return nil
}

func init() {
	business.Register("wac", Init)
}
