package template

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// 如果针对某个id的操作，使用/:id/xxx的形式，尽量避免冲突
// 如果针对多个id的操作，可以使用/multi-xxx的形式，类似multi-delete,然后在body中携带参数
// 如果GET有多个/xxx的，可以在group当中就区分开来，将/:id形式的单独放一个组，其他有具体/xxx的放一个组

// AuthGroup template auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/template")
	{
		g.GET("/", List)
		g.POST("/", Create)
		g.POST("/multi-delete", MultiDelete)
		g.GET("/:id", Get)
		g.DELETE("/:id", Delete)
		g.PUT("/:id", Edit)
	}
}

func init() {
	webserver.RegisterAuthAPI("template", AuthGroup)
}
