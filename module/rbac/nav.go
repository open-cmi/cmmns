package rbac

import (
	"github.com/open-cmi/cmmns/module/license"
	"github.com/open-cmi/gobase/essential/logger"
)

func GetRoleMenus(roleName string) []Menu {
	// 验证license
	// 如果未配置跳过检查，且License检查报错，强制返回未授权菜单配置
	if !gRbacMenuConf.IgnoreLic && license.LicenseCheckError() != nil {
		menus, ok := gRbacMenuConf.Roles["nolic"]
		if !ok {
			logger.Debugf("no menu configuration found for nolic\n")
			return []Menu{}
		}
		return menus
	}

	// 直接从内存配置中获取对应角色的菜单
	menus, ok := gRbacMenuConf.Roles[roleName]
	if !ok {
		logger.Debugf("no menu configuration found for role: %s", roleName)
		return []Menu{}
	}
	return menus
}
