package rbac

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/pkg/goparam"
)

// RouterGroup is a function to api group router
type RouterGroup func(e *gin.Engine)

var mustAuthGroup map[string]string = make(map[string]string)
var optionAuthGroup map[string]string = make(map[string]string)
var unauthGroup map[string]string = make(map[string]string)

// UnauthInit no auth router init
func UnauthInit(e *gin.Engine) {
	for mod, groupPath := range unauthGroup {
		g := e.Group(groupPath)
		{
			modPath, found := unauthAPIPath[mod]
			if !found {
				continue
			}

			for _, r := range modPath {
				switch r.Method {
				case "POST":
					g.POST(r.Path, r.Callback)
				case "GET":
					g.GET(r.Path, r.Callback)
				case "DELETE":
					g.DELETE(r.Path, r.Callback)
				case "PUT":
					g.PUT(r.Path, r.Callback)
				}
			}
		}
	}
}

// OptionAuthInit auth router init
func OptionAuthInit(e *gin.Engine) {
	for mod, groupPath := range optionAuthGroup {
		g := e.Group(groupPath)
		{
			modPath, found := optionAuthAPIPath[mod]
			if !found {
				continue
			}

			for _, r := range modPath {
				switch r.Method {
				case "POST":
					g.POST(r.Path, r.Callback)
				case "GET":
					g.GET(r.Path, r.Callback)
				case "DELETE":
					g.DELETE(r.Path, r.Callback)
				case "PUT":
					g.PUT(r.Path, r.Callback)
				}
			}
		}
	}
}

// OptionAuthInit auth router init
func MustAuthInit(e *gin.Engine) {
	for mod, groupPath := range mustAuthGroup {
		g := e.Group(groupPath)
		{
			modPath, found := mustAuthAPIPath[mod]
			if !found {
				continue
			}

			for _, r := range modPath {
				switch r.Method {
				case "POST":
					g.POST(r.Path, r.Callback)
				case "GET":
					g.GET(r.Path, r.Callback)
				case "DELETE":
					g.DELETE(r.Path, r.Callback)
				case "PUT":
					g.PUT(r.Path, r.Callback)
				}
			}
		}
	}
}

func RegisterOptionAuthRouter(module string, groupPath string) error {
	_, found := optionAuthGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s auth group api has been registered", module)
		return errors.New(errMsg)
	}
	optionAuthGroup[module] = groupPath
	return nil
}

func RegisterMustAuthRouter(module string, groupPath string) error {
	_, found := mustAuthGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s auth group api has been registered", module)
		return errors.New(errMsg)
	}
	mustAuthGroup[module] = groupPath
	return nil
}

func RegisterUnauthRouter(module string, groupPath string) error {
	_, found := unauthGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s unauth group api has been registered", module)
		return errors.New(errMsg)
	}
	unauthGroup[module] = groupPath
	return nil
}

type API struct {
	Prod     string
	Method   string
	Path     string
	Perm     string
	Callback func(c *gin.Context)
}

var mustAuthAPIPath map[string][]API = make(map[string][]API)
var optionAuthAPIPath map[string][]API = make(map[string][]API)
var unauthAPIPath map[string][]API = make(map[string][]API)

func UnauthAPI(prod, method, path string, handler gin.HandlerFunc) error {
	modPath, found := unauthAPIPath[prod]
	if !found {
		unauthAPIPath[prod] = []API{}
		modPath = unauthAPIPath[prod]
	}

	modPath = append(modPath, API{
		Prod:     prod,
		Method:   method,
		Path:     path,
		Perm:     "",
		Callback: handler,
	})
	unauthAPIPath[prod] = modPath

	return nil
}

// MustAuthAPI registers a must-auth API and enforces permission codes.
// If perms is empty, a default perm code is generated from (prod, method, path).
func MustAuthAPI(prod, method, path string, handler gin.HandlerFunc, perm string) error {
	regMu.Lock()
	regAll = append(regAll, APIPermDef{
		Prod:   prod,
		Method: method,
		Path:   path,
		Perm:   perm,
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
			c.String(http.StatusForbidden, i18n.Sprintf("no permission"))
			return
		}
		handler(c)
	}

	modPath, found := mustAuthAPIPath[prod]
	if !found {
		mustAuthAPIPath[prod] = []API{}
		modPath = mustAuthAPIPath[prod]
	}

	modPath = append(modPath, API{
		Prod:     prod,
		Method:   method,
		Path:     path,
		Callback: wrapped,
	})
	mustAuthAPIPath[prod] = modPath

	return nil
}

// AuthAPI registers an auth API and enforces permission codes.
// If perms is empty, a default perm code is generated from (prod, method, path).
func OptionAuthAPI(prod, method, path string, handler gin.HandlerFunc, perm string) error {
	regMu.Lock()
	regAll = append(regAll, APIPermDef{
		Prod:   prod,
		Method: method,
		Path:   path,
		Perm:   perm,
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
			c.String(http.StatusForbidden, i18n.Sprintf("no permission"))
			return
		}
		handler(c)
	}

	modPath, found := optionAuthAPIPath[prod]
	if !found {
		optionAuthAPIPath[prod] = []API{}
		modPath = optionAuthAPIPath[prod]
	}

	modPath = append(modPath, API{
		Prod:     prod,
		Method:   method,
		Path:     path,
		Callback: wrapped,
	})
	optionAuthAPIPath[prod] = modPath

	return nil
}
