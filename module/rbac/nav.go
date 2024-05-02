package rbac

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/system/prod"
)

func GetNavMenu(roleName string) []prod.Menu {
	// 普通用户获取授权菜单
	if roleName == "admin" {
		menu := prod.GetProdNav()
		return menu
	}
	var menu []prod.Menu = []prod.Menu{}
	role := GetByName(roleName)
	if role == nil {
		logger.Errorf("role not exist")
		return menu
	}
	if role.Permisions == "*" {
		return prod.GetProdNav()
	}
	err := json.Unmarshal([]byte(role.Permisions), &menu)
	if err != nil {
		return menu
	}
	return menu
}
