package rbac

import (
	"errors"
	"strings"
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

// AllowRole checks whether the role has at least one of the specified permission codes (OR relation).
// Supports variadic perms as well as comma-separated strings (e.g. "rbac:view,user:view").
// Special-case: role "admin" or empty perm list is always allowed.
func AllowRole(role string, perms ...string) bool {
	if role == "" {
		return false
	}

	if role == "admin" || len(perms) == 0 {
		return true
	}

	// 展开所有可能的权限码（支持逗号分隔与变长参数）
	var targetPerms []string
	for _, p := range perms {
		if p == "" {
			return true // 包含空权限要求，视为公开放行
		}
		// 支持逗号分隔格式："permA,permB"
		parts := strings.Split(p, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				targetPerms = append(targetPerms, trimmed)
			}
		}
	}

	if len(targetPerms) == 0 {
		return true
	}

	userPerms, err := getRolePerms(role)
	if err != nil {
		return false
	}

	// 只要当前角色拥有目标权限列表中的任意一个，即允许访问（OR 逻辑）
	for _, up := range userPerms {
		if up == "*" {
			return true
		}
		for _, tp := range targetPerms {
			if up == tp {
				return true
			}
		}
	}
	return false
}

// GetRolePerms returns all permissions for a role.
// Role "admin" returns ["*"] indicating full permissions.
func GetRolePerms(role string) ([]string, error) {
	if role == "admin" {
		return []string{"*"}, nil
	}
	return getRolePerms(role)
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
