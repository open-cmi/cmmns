package wac

import (
	"net/netip"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"

	"github.com/open-cmi/cmmns/module/nginxconf"
)

var globalSeq int

func GetWAC() Model {
	m := Get()
	if m == nil {
		m = New()
		m.Mode = "blacklist"
	}
	return *m
}

func SetWAC(req *SetRequest) error {
	m := Get()
	if m == nil {
		m = New()
	}

	globalSeq++
	m.Enable = req.Enable
	m.Mode = req.Mode
	m.RawWhitelist = req.RawWhitelist
	m.RawBlacklist = req.RawBlacklist
	m.Seq = globalSeq

	// 保存之前先应用nginx配置，nginx应用成功再保存
	var err error
	if req.Enable {
		err = nginxconf.ApplyNginxAccessControl(m.Mode, m.RawBlacklist, m.RawWhitelist)
	} else {
		err = nginxconf.CancelNginxAccessControl()
	}
	if err != nil {
		logger.Errorf("wac apply nginx access control failed: %s\n", err.Error())
		return err
	}

	err = m.Save()
	if err != nil {
		logger.Errorf("wac save failed: %s\n", err.Error())
		return err
	}

	return err
}

func IPIsInList(rawList string, ip string) bool {
	ipaddr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}

	lines := strings.Split(rawList, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "\t ")
		if line == "" {
			continue
		}
		if strings.Contains(line, "/") {
			// prefix
			p, err := netip.ParsePrefix(line)
			if err != nil {
				continue
			}
			if p.Contains(ipaddr) {
				return true
			}
		} else {
			// addr
			a, err := netip.ParseAddr(line)
			if err != nil {
				continue
			}
			if a.Compare(ipaddr) == 0 {
				return true
			}
		}
	}
	return false
}

func CheckPermit(ip string) bool {

	if globalModel == nil || globalModel.Seq != globalSeq {
		globalModel = Get()
	}
	// 如果无配置，默认是允许的
	if globalModel == nil {
		return true
	}

	// 未开启，默认放行
	if !globalModel.Enable {
		return true
	}

	if globalModel.Mode == "blacklist" {
		// 默认通过
		return !IPIsInList(globalModel.RawBlacklist, ip)
	}
	// 如果是白名单，则在列表中放行，不在则deny
	return IPIsInList(globalModel.RawWhitelist, ip)
}
