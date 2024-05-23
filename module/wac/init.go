package wac

import (
	"github.com/open-cmi/cmmns/module/nginxconf"
	"github.com/open-cmi/cmmns/service/initial"
)

func Init() error {
	m := Get()
	if m != nil {
		if m.Enable {
			return nginxconf.ApplyNginxAccessControl(m.Mode, m.RawBlacklist, m.RawWhitelist)
		}
	}
	return nil
}

func init() {
	initial.Register("wac", initial.DefaultPriority, Init)
}
