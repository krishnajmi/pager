package server

import "sync"

var (
	permissionCache = struct {
		sync.RWMutex
		m map[string][]string // key: "METHOD:path", value: required permissions
	}{m: make(map[string][]string)}
)

// CachePermissions stores required permissions for an API endpoint
// Path should be the relative path (without router group prefix)
func CachePermissions(method, path string, permissions []string) {
	key := method + ":" + path
	permissionCache.Lock()
	defer permissionCache.Unlock()
	permissionCache.m[key] = permissions
}

// GetCachedPermissions retrieves cached permissions for an API endpoint
// fullPath should include the router group prefix
func GetCachedPermissions(method, fullPath string) ([]string, bool) {
	permissionCache.RLock()
	defer permissionCache.RUnlock()

	key := method + ":" + fullPath
	perms, exists := permissionCache.m[key]
	return perms, exists
}
