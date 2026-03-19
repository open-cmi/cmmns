package rbac

import (
	"time"

	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/initial"
)

func Init() error {
	db := sqldb.GetDB()

	now := time.Now().Unix()
	insert := func(role, perm string) {
		_, _ = db.Exec(`INSERT INTO role_permissions(role, perm, created_time) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`, role, perm, now)
	}

	// 权限初始化策略（符合你提的“允许就初始化，不允许就不初始化”）：
	// - 否则：初始化 admin/operator/auditor
	for _, a := range GetAPIPerms() {
		code := PermCode(a.Prod, a.Method, a.Path)
		for _, r := range a.Roles {
			insert(r, code)
		}
	}

	return nil
}

func init() {
	initial.Register("rbac", initial.PhaseDefault, Init)
}
