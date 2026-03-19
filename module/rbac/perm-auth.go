package rbac

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/essential/webserver"
	"github.com/open-cmi/gobase/pkg/goparam"
)

type APIPermDef struct {
	Prod   string
	Method string
	Path   string
	Roles  []string
}

var (
	regMu  sync.Mutex
	regAll []APIPermDef
)

// PermCode generates a stable permission code for an API.
// Convention: "<prod>:<method>:<path>", with method lower-cased.
func PermCode(prod, method, path string) string {
	return prod + ":" + strings.ToLower(strings.TrimSpace(method)) + ":" + strings.TrimSpace(path)
}

func UnauthAPI(prod, method, path string, handler gin.HandlerFunc) error {
	return webserver.RegisterUnauthAPI(prod, method, path, handler)
}

// MustAuthAPI registers a must-auth API and enforces permission codes.
// If perms is empty, a default perm code is generated from (prod, method, path).
func MustAuthAPI(prod, method, path string, handler gin.HandlerFunc, initRoles []string) error {
	perm := PermCode(prod, method, path)

	regMu.Lock()
	regAll = append(regAll, APIPermDef{
		Prod:   prod,
		Method: method,
		Path:   path,
		Roles:  initRoles,
	})
	regMu.Unlock()

	wrapped := func(c *gin.Context) {
		user := goparam.GetUser(c)
		if user == nil {
			c.String(http.StatusUnauthorized, "authenticate is required")
			return
		}

		role, _ := user["role"].(string)
		allowed := AllowRole(role, perm)
		if !allowed {
			c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": i18n.Sprintf("no permission")})
			return
		}
		handler(c)
	}

	return webserver.RegisterMustAuthAPI(prod, method, path, wrapped)
}

// AuthAPI registers an auth API and enforces permission codes.
// If perms is empty, a default perm code is generated from (prod, method, path).
func OptionAuthAPI(prod, method, path string, handler gin.HandlerFunc, initRoles []string) error {
	perm := PermCode(prod, method, path)

	regMu.Lock()
	regAll = append(regAll, APIPermDef{
		Prod:   prod,
		Method: method,
		Path:   path,
		Roles:  initRoles,
	})
	regMu.Unlock()

	wrapped := func(c *gin.Context) {
		user := goparam.GetUser(c)
		if user == nil {
			c.String(http.StatusUnauthorized, "authenticate is required")
			return
		}

		role, _ := user["role"].(string)
		allowed := AllowRole(role, perm)
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"ret": -1, "msg": i18n.Sprintf("no permission")})
			return
		}
		handler(c)
	}

	return webserver.RegisterOptionAuthAPI(prod, method, path, wrapped)
}

func GetAPIPerms() []APIPermDef {
	regMu.Lock()
	out := make([]APIPermDef, len(regAll))
	copy(out, regAll)
	regMu.Unlock()
	return out
}

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
