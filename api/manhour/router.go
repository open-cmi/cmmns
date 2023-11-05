package manhour

import (
	"github.com/open-cmi/cmmns/service/webserver"
)

func init() {
	webserver.RegisterAuthRouter("manhour", "/api/common/v3/manhour")
	webserver.RegisterAuthAPI("manhour", "GET", "/", List)
	webserver.RegisterAuthAPI("manhour", "POST", "/", Create)
	webserver.RegisterAuthAPI("manhour", "POST", "/multi-delete", MultiDelete)
	webserver.RegisterAuthAPI("manhour", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("manhour", "DELETE", "/:id", Delete)
	webserver.RegisterAuthAPI("manhour", "PUT", "/:id", Edit)
}
