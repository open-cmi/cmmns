package wac

import (
	"github.com/open-cmi/cmmns/module/nginxconf"
	"github.com/open-cmi/cmmns/module/wac/blacklist"
	"github.com/open-cmi/cmmns/module/wac/whitelist"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/initial"
)

func Init() error {
	m := Get()
	if m != nil {
		if m.Enable {
			_, lists, err := blacklist.List(nil)
			if err != nil {
				return err
			}
			var blacklists []string = []string{}
			for _, b := range lists {
				blacklists = append(blacklists, b.Address)
			}
			err = nginxconf.ApplyNginxBlackConf(blacklists)
			if err != nil {
				logger.Errorf("apply nginx blacklist conf failed: %s\n", err.Error())
			}

			// 白名单
			_, wlists, err := whitelist.List(nil)
			if err != nil {
				return err
			}
			var whitelists []string = []string{}
			for _, w := range wlists {
				whitelists = append(whitelists, w.Address)
			}
			err = nginxconf.ApplyNginxWhiteConf(whitelists)
			if err != nil {
				logger.Errorf("apply nginx whitelist conf failed: %s\n", err.Error())
			}
			return nginxconf.ApplyNginxAccessControl(m.Mode)
		}
	}
	return nil
}

func init() {
	initial.Register("wac", initial.PhaseDefault, Init)
}
