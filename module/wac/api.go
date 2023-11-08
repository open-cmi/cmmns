package wac

import (
	"net/netip"

	"github.com/open-cmi/cmmns/essential/logger"
)

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
	m.Mode = req.Mode
	m.RawPermitAddress = req.RawPermitAddress
	m.RawDenyAddress = req.RawDenyAddress
	err := m.Save()
	if err != nil {
		logger.Errorf("wac save failed: %s\n", err.Error())
		return err
	}
	gWebAccessControl = m.ConvertoWAC()
	return nil
}

func CheckPermit(ip string) bool {
	var permit bool
	if gWebAccessControl.Mode == "blacklist" {
		// 默认通过
		permit = true
		ipaddr, err := netip.ParseAddr(ip)
		if err != nil {
			return false
		}
		for _, addr := range gWebAccessControl.DenyAddrs {
			if addr.Compare(ipaddr) == 0 {
				return false
			}
		}
		for _, prefix := range gWebAccessControl.DenyPrefixs {
			if prefix.Contains(ipaddr) {
				return false
			}
		}
	} else {
		// whitelist, 默认允许
		permit = false
		// 默认通过
		permit = true
		ipaddr, err := netip.ParseAddr(ip)
		if err != nil {
			return false
		}
		for _, addr := range gWebAccessControl.PermitAddrs {
			if addr.Compare(ipaddr) == 0 {
				return true
			}
		}
		for _, prefix := range gWebAccessControl.PermitPrefixs {
			if prefix.Contains(ipaddr) {
				return true
			}
		}
	}
	return permit
}
