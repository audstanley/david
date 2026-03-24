package auth

import "time"

// AuthResult represents the result of authentication
type AuthResult struct {
	UserID      string
	Role        string
	Permissions []string
	AuthMethod  string
	IP          string
}

// AuthContext is used to pass auth info through middleware
type AuthContext struct {
	Result *AuthResult
}

// APIKey represents an API key
type APIKey struct {
	ID          string
	UserID      string
	Name        string
	Hash        string
	Permissions []string
	Expiry      time.Time
	Created     time.Time
	LastUsed    time.Time
	Disabled    bool
}

// Role-based access control roles
const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleReader = "reader"
	RolePublic = "public"
)

// Permission constants
const (
	// Calendar permissions
	PermCalendarCreate = "calendar.create"
	PermCalendarRead   = "calendar.read"
	PermCalendarUpdate = "calendar.update"
	PermCalendarDelete = "calendar.delete"
	PermCalendarShare  = "calendar.share"

	// Event permissions
	PermEventCreate = "event.create"
	PermEventRead   = "event.read"
	PermEventUpdate = "event.update"
	PermEventDelete = "event.delete"

	// Todo permissions
	PermTodoCreate = "todo.create"
	PermTodoRead   = "todo.read"
	PermTodoUpdate = "todo.update"
	PermTodoDelete = "todo.delete"

	// Journal permissions
	PermJournalCreate = "journal.create"
	PermJournalRead   = "journal.read"
	PermJournalUpdate = "journal.update"
	PermJournalDelete = "journal.delete"

	// User permissions (admin only)
	PermUserCreate = "user.create"
	PermUserRead   = "user.read"
	PermUserUpdate = "user.update"
	PermUserDelete = "user.delete"

	// System permissions (admin only)
	PermSystemConfig = "system.config"
)

// Permission sets for each role
var RolePermissions = map[string][]string{
	RoleAdmin: {
		PermCalendarCreate, PermCalendarRead, PermCalendarUpdate, PermCalendarDelete, PermCalendarShare,
		PermEventCreate, PermEventRead, PermEventUpdate, PermEventDelete,
		PermTodoCreate, PermTodoRead, PermTodoUpdate, PermTodoDelete,
		PermJournalCreate, PermJournalRead, PermJournalUpdate, PermJournalDelete,
		PermUserCreate, PermUserRead, PermUserUpdate, PermUserDelete,
		PermSystemConfig,
	},
	RoleUser: {
		PermCalendarCreate, PermCalendarRead, PermCalendarUpdate, PermCalendarDelete, PermCalendarShare,
		PermEventCreate, PermEventRead, PermEventUpdate, PermEventDelete,
		PermTodoCreate, PermTodoRead, PermTodoUpdate, PermTodoDelete,
		PermJournalCreate, PermJournalRead, PermJournalUpdate, PermJournalDelete,
	},
	RoleReader: {
		PermCalendarRead,
		PermEventRead,
		PermTodoRead,
		PermJournalRead,
	},
	RolePublic: {
		PermCalendarRead,
		PermEventRead,
		PermTodoRead,
		PermJournalRead,
	},
}

// IsPermissionAllowed checks if a role has a specific permission
func IsPermissionAllowed(role, permission string) bool {
	perms, ok := RolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

// HasPermission checks if a list of permissions includes a specific permission
func HasPermission(permissions []string, permission string) bool {
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// Auth methods
const (
	AuthMethodJWT    = "jwt"
	AuthMethodAPIKey = "api_key"
	AuthMethodBasic  = "basic"
)

// JWT claim types
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)
