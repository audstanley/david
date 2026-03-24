# Phase 3: Authentication & Authorization

## Overview
This phase implements multi-layer authentication (JWT, API Key, HTTP Basic Auth) with role-based access control (RBAC). It provides secure access to the API and WebDAV endpoints with proper token management and permission checking.

## Goals
- Implement JWT authentication with refresh tokens
- Implement API key authentication
- Implement HTTP Basic Auth integration
- Create role-based access control (Admin, User, Reader)
- Implement token blacklist for revocation
- Implement rate limiting per IP/user
- Achieve 90%+ unit test coverage

## Tasks

### 1. JWT Authentication
- [ ] Create `app/auth/jwt/token.go` - JWT token management
  - [ ] Generate access token (short-lived, e.g., 24h)
  - [ ] Generate refresh token (long-lived, e.g., 7d)
  - [ ] Verify access token signature and expiry
  - [ ] Verify refresh token signature and expiry
  - [ ] Extract user ID from token
  - [ ] Extract role from token
  - [ ] Support multiple claims (user_id, role, token_type)
- [ ] Create `app/auth/jwt/refresh.go` - Refresh token logic
  - [ ] Exchange refresh token for new access token
  - [ ] Validate refresh token is not revoked
  - [ ] Handle refresh token rotation (optional)
  - [ ] Invalidate old refresh token after use (optional)
- [ ] Create `app/auth/jwt/blacklist.go` - Token revocation
  - [ ] Store revoked token IDs in LevelDB
  - [ ] Check blacklist during token verification
  - [ ] Implement blacklist cleanup (expired entries)
  - [ ] Support forced logout (add to blacklist)
- [ ] Create `app/auth/jwt/types.go` - JWT types

### 2. API Key Authentication
- [ ] Create `app/auth/apikey/generator.go` - API key generation
  - [ ] Generate random API key (cryptographically secure)
  - [ ] Format key with prefix (e.g., "david_")
  - [ ] Generate key hash for storage (never store plain key)
- [ ] Create `app/auth/apikey/verifier.go` - API key validation
  - [ ] Verify API key format
  - [ ] Hash provided key and compare to stored hash
  - [ ] Check API key expiry
  - [ ] Check API key is not disabled
  - [ ] Extract user_id and permissions from key metadata
- [ ] Create `app/auth/apikey/store.go` - API key storage
  - [ ] Store new API key (hash + metadata)
  - [ ] Get API key by ID
  - [ ] List API keys for user
  - [ ] Revoke API key
  - [ ] Update API key metadata
- [ ] Create `app/auth/apikey/types.go` - API key types
  - [ ] KeyID, UserID, Name, Hash
  - [ ] Permissions (read, write, admin)
  - [ ] Expiry, Created, LastUsed

### 3. HTTP Basic Auth
- [ ] Create `app/auth/basic/handler.go` - HTTP Basic Auth handler
  - [ ] Parse Authorization: Basic header
  - [ ] Decode base64 username:password
  - [ ] Validate username format
  - [ ] Delegate password verification to user store
- [ ] Create `app/auth/basic/verifier.go` - Password verification
  - [ ] Get user by username
  - [ ] Verify password hash (bcrypt, argon2, scrypt)
  - [ ] Return user on success, error on failure
  - [ ] Handle account locked status (optional)
  - [ ] Log failed login attempts (for audit)
- [ ] Create `app/auth/basic/types.go` - Basic auth types

### 4. Role-Based Access Control (RBAC)
- [ ] Create `app/auth/rbac/roles.go` - Role definitions
  - [ ] Define Admin role permissions
  - [ ] Define User role permissions
  - [ ] Define Reader role permissions
  - [ ] Define Public role permissions (for public calendars)
  - [ ] Implement role hierarchy (Admin > User > Reader)
- [ ] Create `app/auth/rbac/permissions.go` - Permission definitions
  - [ ] Define granular permissions:
    - calendar.create, calendar.read, calendar.update, calendar.delete
    - calendar.share_grant, calendar.share_revoke
    - event.create, event.read, event.update, event.delete
    - todo.create, todo.read, todo.update, todo.delete
    - journal.create, journal.read, journal.update, journal.delete
    - user.create, user.read, user.update, user.delete (admin only)
    - system.config (admin only)
- [ ] Create `app/auth/rbac/checker.go` - Access control decisions
  - [ ] Check if user has permission for action
  - [ ] Check calendar ownership
  - [ ] Check calendar share permissions
  - [ ] Check public calendar access
  - [ ] Return allow/deny with reason
- [ ] Create `app/auth/rbac/middleware.go` - RBAC middleware
  - [ ] Middleware to check permissions on API routes
  - [ ] Skip check for public endpoints
  - [ ] Return 403 Forbidden if access denied
  - [ ] Include permission error in response

### 5. Rate Limiting
- [ ] Create `app/auth/ratelimit/limiter.go` - Rate limiter
  - [ ] Implement token bucket algorithm
  - [ ] Track requests per IP address
  - [ ] Track requests per user ID (if authenticated)
  - [ ] Configurable limits (default 100 req/min)
  - [ ] Admin override for specific users/calendars
- [ ] Create `app/auth/ratelimit/store.go` - Rate limit storage
  - [ ] Store rate limit counters in LevelDB
  - [ ] Implement counter cleanup (expire old entries)
  - [ ] Support sliding window
- [ ] Create `app/auth/ratelimit/middleware.go` - Rate limit middleware
  - [ ] Middleware to check rate limit
  - [ ] Return 429 Too Many Requests if exceeded
  - [ ] Include rate limit headers (X-RateLimit-Limit, X-RateLimit-Remaining)
  - [ ] Log rate limit violations (for audit)
- [ ] Create `app/auth/ratelimit/admin.go` - Admin rate limit management
  - [ ] List rate limit rules
  - [ ] Create rate limit rule (IP, user, calendar)
  - [ ] Update rate limit rule
  - [ ] Delete rate limit rule

### 6. Auth Manager (Unified Interface)
- [ ] Create `app/auth/manager.go` - Multi-auth manager
  - [ ] Detect auth method from request (JWT, API key, Basic)
  - [ ] Delegate to appropriate auth handler
  - [ ] Normalize auth result (user_id, role, permissions)
  - [ ] Handle auth failures consistently
- [ ] Create `app/auth/types.go` - Unified auth types
  - [ ] AuthResult struct (UserID, Role, Permissions, AuthMethod)
  - [ ] AuthContext struct (for passing through context)

### 7. Middleware Stack
- [ ] Create `app/api/middleware/auth.go` - Auth middleware
  - [ ] Extract auth from request
  - [ ] Verify auth
  - [ ] Add AuthResult to context
  - [ ] Return 401 Unauthorized if failed
- [ ] Create `app/api/middleware/rbac.go` - RBAC middleware
  - [ ] Check permissions based on route
  - [ ] Return 403 Forbidden if denied
- [ ] Create `app/api/middleware/ratelimit.go` - Rate limit middleware
  - [ ] Check rate limit
  - [ ] Return 429 if exceeded
  - [ ] Add rate limit headers
- [ ] Create `app/api/middleware/cors.go` - CORS middleware
  - [ ] Handle CORS preflight
  - [ ] Set CORS headers based on config
- [ ] Create `app/api/middleware/logging.go` - Request logging
  - [ ] Log all requests (method, path, status, duration)
  - [ ] Include user ID if authenticated
  - [ ] Exclude health check endpoints
- [ ] Create `app/api/middleware/recovery.go` - Panic recovery
  - [ ] Recover from panics
  - [ ] Return 500 Internal Server Error
  - [ ] Log panic details

### 8. Authentication Endpoints (API)
- [ ] Create `app/api/handlers/auth.go` - Auth handlers
  - [ ] POST /api/auth/login (Basic Auth → return JWT + refresh)
  - [ ] POST /api/auth/refresh (Exchange refresh token for new access)
  - [ ] POST /api/auth/logout (Add token to blacklist)
  - [ ] GET /api/auth/verify (Verify token, return user info)
  - [ ] POST /api/auth/api-key (Generate new API key - User role)
  - [ ] DELETE /api/auth/api-key/{id} (Revoke API key)
  - [ ] GET /api/auth/api-key (List API keys)

### 9. Unit Tests
- [ ] Create `app/auth/jwt/token_test.go`
  - [ ] Test token generation
  - [ ] Test token verification
  - [ ] Test expiry handling
  - [ ] Test invalid signatures
- [ ] Create `app/auth/jwt/refresh_test.go`
  - [ ] Test refresh token exchange
  - [ ] Test expired refresh tokens
  - [ ] Test revoked refresh tokens
- [ ] Create `app/auth/jwt/blacklist_test.go`
  - [ ] Test token blacklist operations
  - [ ] Test blacklist lookup
- [ ] Create `app/auth/apikey/generator_test.go`
  - [ ] Test key generation (randomness, format)
  - [ ] Test key hashing
- [ ] Create `app/auth/apikey/verifier_test.go`
  - [ ] Test key verification
  - [ ] Test expired keys
  - [ ] Test disabled keys
- [ ] Create `app/auth/apikey/store_test.go`
- [ ] Create `app/auth/basic/handler_test.go`
  - [ ] Test Basic Auth header parsing
  - [ ] Test password verification
  - [ ] Test invalid credentials
- [ ] Create `app/auth/rbac/roles_test.go`
  - [ ] Test role definitions
  - [ ] Test role hierarchy
- [ ] Create `app/auth/rbac/permissions_test.go`
  - [ ] Test permission checks
  - [ ] Test permission combinations
- [ ] Create `app/auth/rbac/checker_test.go`
  - [ ] Test access control decisions
  - [ ] Test ownership checks
  - [ ] Test share permission checks
- [ ] Create `app/auth/ratelimit/limiter_test.go`
  - [ ] Test token bucket algorithm
  - [ ] Test rate limit enforcement
  - [ ] Test admin overrides

## Data Models

### JWT Claims
```go
type Claims struct {
    UserID   string `json:"user_id"`
    Role     string `json:"role"`
    TokenType string `json:"token_type"`  // access, refresh
    JWTID    string `json:"jti"`         // JWT ID for blacklist
    ExpiresAt int64 `json:"exp"`
}
```

### API Key
```go
type APIKey struct {
    ID        string
    UserID    string
    Name      string
    Hash      string  // Hashed key (never store plain)
    Permissions []string  // List of permissions
    Expiry    time.Time
    Created   time.Time
    LastUsed  time.Time
    Disabled  bool
}
```

### Rate Limit Rule
```go
type RateLimitRule struct {
    ID         string
    TargetType string  // ip, user, calendar
    TargetID   string
    Limit      int     // requests per window
    Window     string  // e.g., "1m", "1h"
    Created    time.Time
    CreatedBy  string
}
```

### Auth Result
```go
type AuthResult struct {
    UserID      string
    Role        string
    Permissions []string
    AuthMethod  string  // jwt, api_key, basic
    IP          string
}
```

## Success Criteria
- [ ] JWT authentication works (generate, verify, refresh)
- [ ] JWT blacklist works (revocation)
- [ ] API key authentication works (generate, verify, revoke)
- [ ] HTTP Basic Auth works with password verification
- [ ] Multi-auth manager correctly detects and handles auth methods
- [ ] RBAC roles are correctly defined (Admin, User, Reader, Public)
- [ ] Permission checks work correctly
- [ ] Ownership checks work correctly
- [ ] Share permission checks work correctly
- [ ] Rate limiting works (per IP, per user, admin overrides)
- [ ] Rate limit middleware adds correct headers
- [ ] 90%+ unit test coverage
- [ ] All tests pass with `go test -race ./app/auth/...`

## Notes
- JWT secret must be set via environment variable (never in config file)
- Refresh tokens should be rotated (optional but recommended)
- API keys should be hashed before storage (like passwords)
- Rate limiting should be configurable per endpoint (optional enhancement)
- Auth failures should not leak information (generic error messages)

## Next Phase
Move to **Phase 4: REST API (Gorilla Mux)** when all items above are complete.
