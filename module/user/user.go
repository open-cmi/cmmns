package user

import (
	"github.com/open-cmi/gobase/essential/rdb"
)

func init() {
	rdb.Register("user", rdb.RDBUser)
}
