package assist

import (
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/assist/router"
)

func init() {
	api.RegisterAuthAPI("assist", router.AuthGroup)
}
