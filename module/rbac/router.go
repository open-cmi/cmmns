package rbac

import (
	"github.com/open-cmi/gobase/essential/webserver"
)

// 以下路由注册函数直接委托 webserver 的对应实现，便于统一从 rbac 包入口调用。
// 若 gobase 提供 RegisterAuthRouter、RegisterOptionAuthRouter 等，可在此追加封装后改为 rbac 调用。
func RegisterOptionAuthRouter(module, groupPath string) error {
	return webserver.RegisterOptionAuthRouter(module, groupPath)
}

// RegisterMustAuthRouter 注册必须认证的路由组，直接委托 webserver。
func RegisterMustAuthRouter(module, groupPath string) error {
	return webserver.RegisterMustAuthRouter(module, groupPath)
}

// RegisterUnauthRouter 注册无需认证的路由组，直接委托 webserver。
func RegisterUnauthRouter(module, groupPath string) error {
	return webserver.RegisterUnauthRouter(module, groupPath)
}
