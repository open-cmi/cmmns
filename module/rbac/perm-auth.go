package rbac

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/gobase/essential/i18n"
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
	return RegisterUnauthAPI(prod, method, path, handler)
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

	return RegisterMustAuthAPI(prod, method, path, wrapped)
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

	return RegisterOptionAuthAPI(prod, method, path, wrapped)
}

func GetAPIPerms() []APIPermDef {
	regMu.Lock()
	out := make([]APIPermDef, len(regAll))
	copy(out, regAll)
	regMu.Unlock()
	return out
}
