package rbac

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/license"
)

func GetRoleMenus(roleName string) []Menu {
	// 验证license
	if license.LicenseCheckError() != nil {
		return GetMenusWhileNoLicense()
	}

	role := GetByName(roleName)
	if role == nil {
		logger.Errorf("role is not existing")
		return []Menu{}
	}

	if role.Permisions == "*" {
		menus, ok := gRbacMenus.Roles[roleName]
		if ok {
			return menus
		}
		return []Menu{}
	}

	var menus []Menu = []Menu{}
	err := json.Unmarshal([]byte(role.Permisions), &menus)
	if err != nil {
		logger.Errorf("get role menus unmarshal failed: %s\n", err.Error())
	}
	return menus
}

func GetMenusWhileNoLicense() []Menu {
	return gRbacMenus.NoLic
}
