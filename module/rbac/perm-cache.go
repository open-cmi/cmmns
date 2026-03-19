package rbac

import (
	"errors"
	"sync"

	"github.com/open-cmi/gobase/essential/sqldb"
)

type cachedPerms struct {
	perms      []string
	loadedRole string
}

var (
	permCacheMu sync.Mutex
	permCache   = map[string]cachedPerms{}
)

// AllowRole checks whether the role has all required permission codes.
// Special-case: role "admin" is always allowed.
func AllowRole(role string, perm string) bool {
	if role == "" {
		return false
	}

	if perm == "" {
		return true
	}

	perms, err := getRolePerms(role)
	if err != nil {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return true
}

func getRolePerms(role string) ([]string, error) {
	permCacheMu.Lock()
	c, ok := permCache[role]
	if ok {
		permCacheMu.Unlock()
		return c.perms, nil
	}
	permCacheMu.Unlock()

	db := sqldb.GetDB()
	if db == nil {
		return []string{}, errors.New("db is not available")
	}

	rows, err := db.Queryx(`select perm from role_permissions where role=$1`, role)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	perms := []string{}
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return []string{}, err
		}
		if perm != "" {
			perms = append(perms, perm)
		}
	}

	out := cachedPerms{
		perms:      perms,
		loadedRole: role,
	}

	permCacheMu.Lock()
	permCache[role] = out
	permCacheMu.Unlock()

	return perms, nil
}

// InvalidateRolePermCache clears cache for a role (call after updates).
func InvalidateRolePermCache(role string) {
	permCacheMu.Lock()
	delete(permCache, role)
	permCacheMu.Unlock()
}
