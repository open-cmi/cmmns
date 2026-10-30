package rbac

import (
	"sync"
)

type APIPermDef struct {
	Prod   string
	Method string
	Path   string
	Perm   string
}

var (
	regMu  sync.Mutex
	regAll []APIPermDef
)
