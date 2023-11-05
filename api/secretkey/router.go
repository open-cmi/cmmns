package secretkey

import (
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("secretkey", "/api/secret-key/v1")
	webserver.RegisterAuthAPI("secretkey", "GET", "/", List)
	webserver.RegisterAuthAPI("secretkey", "GET", "/name-list/", NameList)
	webserver.RegisterAuthAPI("secretkey", "POST", "/", Create)
	webserver.RegisterAuthAPI("secretkey", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("secretkey", "DELETE", "/:id", Delete)
	webserver.RegisterAuthAPI("secretkey", "PUT", "/:id", Edit)
}
