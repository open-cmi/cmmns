package rbac

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/license"
	"github.com/open-cmi/cmmns/module/system/prod"
)

func GetNavMenu(roleName string) []prod.Menu {
	// 验证license
	if license.LicenseCheckError() != nil {
		return prod.GetRequireNav()
	}

	// 普通用户获取授权菜单
	if roleName == "admin" {
		menu := prod.GetNav()
		return menu
	}
	var menu []prod.Menu = []prod.Menu{}
	role := GetByName(roleName)
	if role == nil {
		logger.Errorf("role not exist")
		return menu
	}
	if role.Permisions == "*" {
		return prod.GetNav()
	}
	err := json.Unmarshal([]byte(role.Permisions), &menu)
	if err != nil {
		return menu
	}
	return menu
}
