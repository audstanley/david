package rbac

import (
	"github.com/audstanley/david/app/auth/common"
)

// RoleChecker checks role-based access control
type RoleChecker struct {
	permissions map[string][]string
}

// NewRoleChecker creates a new role checker
func NewRoleChecker() *RoleChecker {
	return &RoleChecker{
		permissions: common.RolePermissions,
	}
}

// HasRole checks if a user has a specific role
func (c *RoleChecker) HasRole(userRole, requiredRole string) bool {
	// Role hierarchy: admin > user > reader > public
	switch requiredRole {
	case common.RolePublic:
		return true // Everyone has public access
	case common.RoleReader:
		return userRole == common.RoleReader || userRole == common.RoleUser || userRole == common.RoleAdmin
	case common.RoleUser:
		return userRole == common.RoleUser || userRole == common.RoleAdmin
	case common.RoleAdmin:
		return userRole == common.RoleAdmin
	default:
		return userRole == requiredRole
	}
}

// HasPermission checks if a user has a specific permission
func (c *RoleChecker) HasPermission(userRole, permission string) bool {
	return common.IsPermissionAllowed(userRole, permission)
}

// HasPermissions checks if a user has all required permissions
func (c *RoleChecker) HasPermissions(userRole string, permissions []string) bool {
	for _, perm := range permissions {
		if !c.HasPermission(userRole, perm) {
			return false
		}
	}
	return true
}

// CheckAccess determines if access should be granted
type AccessResult struct {
	Allowed bool
	Reason  string
}

// CheckAccess checks if a user can access a resource
func (c *RoleChecker) CheckAccess(userRole, requiredRole string, userPermissions []string, requiredPermissions []string) *AccessResult {
	// Check role
	if !c.HasRole(userRole, requiredRole) {
		return &AccessResult{
			Allowed: false,
			Reason:  "insufficient role",
		}
	}

	// Check permissions
	if len(requiredPermissions) > 0 {
		if !c.HasPermissions(userRole, requiredPermissions) {
			return &AccessResult{
				Allowed: false,
				Reason:  "missing permissions",
			}
		}
	}

	return &AccessResult{
		Allowed: true,
		Reason:  "access granted",
	}
}

// GetRolePermissions returns the permissions for a role
func (c *RoleChecker) GetRolePermissions(role string) []string {
	perms, ok := c.permissions[role]
	if !ok {
		return []string{}
	}
	return perms
}
