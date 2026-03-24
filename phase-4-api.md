# Phase 4: REST API (Gorilla Mux)

## Overview
This phase implements the REST API using Gorilla Mux. The API provides programmatic access to calendars, events, todos, journals, and admin functions. It supports multiple authentication methods and role-based access control.

## Goals
- Implement REST API with Gorilla Mux
- Create comprehensive API endpoints for all resources
- Implement request/response models with validation
- Implement error handling with standardized responses
- Achieve 85%+ test coverage (integration tests)

## Tasks

### 1. API Router Setup
- [ ] Create `app/api/router.go` - Main router
  - [ ] Initialize Gorilla Mux router
  - [ ] Set up route groups (auth, users, calendars, events, etc.)
  - [ ] Apply middleware stack (auth, rbac, rate limit, CORS, logging, recovery)
  - [ ] Set up health check endpoint
  - [ ] Set up documentation endpoint
  - [ ] Configure route parameters
- [ ] Create `app/api/middleware/auth.go` - Auth middleware (from Phase 3)
  - [ ] Extract and verify authentication
  - [ ] Add AuthResult to context
  - [ ] Return 401 on failure
- [ ] Create `app/api/middleware/rbac.go` - RBAC middleware (from Phase 3)
  - [ ] Check permissions for route
  - [ ] Return 403 on failure
- [ ] Create `app/api/middleware/ratelimit.go` - Rate limit middleware (from Phase 3)
- [ ] Create `app/api/middleware/cors.go` - CORS middleware
  - [ ] Handle OPTIONS preflight
  - [ ] Set Access-Control-Allow-Origin
  - [ ] Set Access-Control-Allow-Methods
  - [ ] Set Access-Control-Allow-Headers
  - [ ] Set Access-Control-Allow-Credentials
- [ ] Create `app/api/middleware/logging.go` - Request logging
- [ ] Create `app/api/middleware/recovery.go` - Panic recovery

### 2. Response Models
- [ ] Create `app/api/models/response.go` - Standard response format
  - [ ] Success response: `{success: true, data: ..., meta: ...}`
  - [ ] Error response: `{success: false, error: {code, message, details}}`
  - [ ] Pagination meta: `{total, page, per_page, has_more}`
  - [ ] Implement JSON marshaling
- [ ] Create `app/api/models/request.go` - Request DTOs
  - [ ] CreateRequest (calendar, event, todo, journal)
  - [ ] UpdateRequest (calendar, event, todo, journal)
  - [ ] ShareRequest (grant/revoke)
  - [ ] QueryParams (date range, filters)
  - [ ] AuthLoginRequest (username, password)
  - [ ] AuthRefreshRequest (refresh_token)
  - [ ] APIKeyRequest (name, permissions)
- [ ] Create `app/api/models/validation.go` - Request validation
  - [ ] Validate required fields
  - [ ] Validate field formats (email, UID, date-time)
  - [ ] Validate length constraints
  - [ ] Return validation errors

### 3. Error Handling
- [ ] Create `app/api/errors.go` - API error types
  - [ ] NotFoundError (404)
  - [ ] UnauthorizedError (401)
  - [ ] ForbiddenError (403)
  - [ ] ValidationError (400)
  - [ ] ConflictError (409)
  - [ ] RateLimitError (429)
  - [ ] InternalError (500)
  - [ ] Implement error to HTTP status mapping
  - [ ] Implement error JSON serialization
- [ ] Create `app/api/handlers/errors.go` - Error handler middleware
  - [ ] Catch and format all errors
  - [ ] Return consistent error responses
  - [ ] Log errors with context

### 4. Authentication Handlers
- [ ] Create `app/api/handlers/auth.go` - Auth endpoints
  - [ ] POST /api/auth/login
    - [ ] Accept Basic Auth credentials
    - [ ] Verify credentials
    - [ ] Generate access + refresh tokens
    - [ ] Return tokens + user info
  - [ ] POST /api/auth/refresh
    - [ ] Accept refresh token
    - [ ] Verify refresh token
    - [ ] Generate new access token
    - [ ] Return access token
  - [ ] POST /api/auth/logout
    - [ ] Accept access token
    - [ ] Add to blacklist
    - [ ] Return success
  - [ ] GET /api/auth/verify
    - [ ] Verify token
    - [ ] Return user info + permissions
  - [ ] POST /api/auth/api-key
    - [ ] Create new API key
    - [ ] Return API key (only once!)
  - [ ] GET /api/auth/api-key
    - [ ] List user's API keys
  - [ ] DELETE /api/auth/api-key/{id}
    - [ ] Revoke API key
- [ ] Create `app/api/handlers/auth_test.go` - Auth handler tests
  - [ ] Test login flow
  - [ ] Test refresh flow
  - [ ] Test logout (blacklist)
  - [ ] Test API key management

### 5. User Handlers (Admin)
- [ ] Create `app/api/handlers/users.go` - User endpoints
  - [ ] GET /api/users
    - [ ] List all users (admin only)
    - [ ] Support pagination
    - [ ] Support filters (role, created date)
  - [ ] GET /api/users/{id}
    - [ ] Get user by ID
  - [ ] POST /api/users
    - [ ] Create new user
    - [ ] Require admin role
    - [ ] Accept password (hash server-side)
    - [ ] Accept role assignment
  - [ ] PUT /api/users/{id}
    - [ ] Update user
    - [ ] Require admin role or user updating own profile
    - [ ] Accept name, email updates
    - [ ] Accept role changes (admin only)
  - [ ] DELETE /api/users/{id}
    - [ ] Delete user
    - [ ] Require admin role
    - [ ] Handle user's calendars (reassign or delete)
- [ ] Create `app/api/handlers/users_test.go`

### 6. Calendar Handlers
- [ ] Create `app/api/handlers/calendars.go` - Calendar endpoints
  - [ ] GET /api/calendars
    - [ ] List calendars for user
    - [ ] Include shared calendars
    - [ ] Support pagination
    - [ ] Support filters (public, owned)
  - [ ] GET /api/calendars/{uid}
    - [ ] Get calendar by UID
    - [ ] Check access (owner, share, public)
  - [ ] POST /api/calendars
    - [ ] Create new calendar
    - [ ] Accept display_name, description, color
    - [ ] Accept is_public flag
    - [ ] Generate public hash if public
  - [ ] PUT /api/calendars/{uid}
    - [ ] Update calendar
    - [ ] Require write access
  - [ ] DELETE /api/calendars/{uid}
    - [ ] Delete calendar
    - [ ] Require admin access
  - [ ] GET /api/calendars/{uid}/export.ics
    - [ ] Export calendar to ICS
  - [ ] GET /api/calendars/{uid}/stats
    - [ ] Get calendar statistics (event count, storage size)
- [ ] Create `app/api/handlers/calendars_test.go`

### 7. Share Handlers
- [ ] Create `app/api/handlers/shares.go` - Share endpoints
  - [ ] GET /api/calendars/{uid}/shares
    - [ ] List shares for calendar
    - [ ] Require write/admin access
  - [ ] POST /api/calendars/{uid}/shares
    - [ ] Grant access to user
    - [ ] Accept user_id, role (read, write, admin, everyone)
    - [ ] Validate user exists
    - [ ] Check for duplicate share
  - [ ] DELETE /api/calendars/{uid}/shares/{id}
    - [ ] Revoke share
    - [ ] Require write/admin access
  - [ ] POST /api/calendars/{uid}/share/public
    - [ ] Make calendar public
    - [ ] Generate/rotate public hash
  - [ ] DELETE /api/calendars/{uid}/share/public
    - [ ] Remove public access
- [ ] Create `app/api/handlers/shares_test.go`

### 8. Event Handlers
- [ ] Create `app/api/handlers/events.go` - Event endpoints
  - [ ] GET /api/events
    - [ ] List events with filters
    - [ ] Query params: calendar, start, end, status, attendee
    - [ ] Support pagination
    - [ ] Return expanded recurrence instances
  - [ ] GET /api/events/{uid}
    - [ ] Get event by UID
    - [ ] Check calendar access
  - [ ] POST /api/events
    - [ ] Create new event
    - [ ] Accept full event structure (ICS-like)
    - [ ] Generate UID if not provided
    - [ ] Validate required properties
  - [ ] PUT /api/events/{uid}
    - [ ] Update event
    - [ ] Increment sequence number
    - [ ] Require write access
    - [ ] Validate RRULE changes
  - [ ] DELETE /api/events/{uid}
    - [ ] Delete event
    - [ ] Require write access
    - [ ] Handle recurrence instances
  - [ ] GET /api/events/{uid}/export.ics
    - [ ] Export single event to ICS
- [ ] Create `app/api/handlers/events_test.go`
  - [ ] Test date range queries
  - [ ] Test recurrence expansion
  - [ ] Test event CRUD

### 9. Todo Handlers
- [ ] Create `app/api/handlers/todos.go` - Todo endpoints
  - [ ] GET /api/todos
    - [ ] List todos with filters
    - [ ] Query params: calendar, status, priority
  - [ ] GET /api/todos/{uid}
  - [ ] POST /api/todos
  - [ ] PUT /api/todos/{uid}
  - [ ] DELETE /api/todos/{uid}
  - [ ] GET /api/todos/{uid}/export.ics
- [ ] Create `app/api/handlers/todos_test.go`

### 10. Journal Handlers
- [ ] Create `app/api/handlers/journals.go` - Journal endpoints
  - [ ] GET /api/journals
  - [ ] GET /api/journals/{uid}
  - [ ] POST /api/journals
  - [ ] PUT /api/journals/{uid}
  - [ ] DELETE /api/journals/{uid}
  - [ ] GET /api/journals/{uid}/export.ics
- [ ] Create `app/api/handlers/journals_test.go`

### 11. Free/Busy Handlers
- [ ] Create `app/api/handlers/freebusy.go` - Free/Busy endpoints
  - [ ] GET /api/freebusy/{calendar_uid}
    - [ ] Generate free/busy for calendar
    - [ ] Query params: start, end
  - [ ] GET /api/freebusy
    - [ ] Aggregate free/busy across multiple calendars
    - [ ] Query params: calendars (comma-separated), start, end
  - [ ] Return ICS VFREEBUSY format
- [ ] Create `app/api/handlers/freebusy_test.go`

### 12. Export Handlers
- [ ] Create `app/api/handlers/export.go` - Export endpoints
  - [ ] GET /api/export/calendar/{uid}.ics
    - [ ] Export calendar to ICS
  - [ ] GET /api/export/event/{uid}.ics
    - [ ] Export single event to ICS
  - [ ] GET /api/export/range.ics
    - [ ] Export events in date range
    - [ ] Query params: calendar, start, end
  - [ ] GET /api/export/all.ics
    - [ ] Export all user's calendars
- [ ] Create `app/api/handlers/export_test.go`

### 13. Admin Handlers
- [ ] Create `app/api/handlers/admin.go` - Admin endpoints
  - [ ] GET /api/admin/statistics
    - [ ] Get system statistics
    - [ ] Total users, calendars, events, storage size
  - [ ] GET /api/admin/statistics/users
    - [ ] Get user statistics (calendars owned, events created)
  - [ ] GET /api/admin/statistics/calendars
    - [ ] Get calendar statistics (event counts, storage)
  - [ ] GET /api/admin/audit
    - [ ] Get audit log
    - [ ] Query params: user, entity, date range
  - [ ] GET /api/admin/rate-limits
    - [ ] List rate limit rules
  - [ ] POST /api/admin/rate-limits
    - [ ] Create rate limit rule
  - [ ] DELETE /api/admin/rate-limits/{id}
    - [ ] Delete rate limit rule
- [ ] Create `app/api/handlers/admin_test.go`

### 14. Health & Status Handlers
- [ ] Create `app/api/handlers/health.go` - Health endpoints
  - [ ] GET /api/health
    - [ ] Return system health status
    - [ ] Check database connectivity
    - [ ] Return status + timestamp
  - [ ] GET /api/status
    - [ ] Return detailed status
    - [ ] Database stats
    - [ ] API stats (request count, errors)
    - [ ] Uptime
- [ ] Create `app/api/handlers/health_test.go`

### 15. Integration Tests
- [ ] Create `app/api/integration_test.go`
  - [ ] Test full auth flow (login → API call → logout)
  - [ ] Test RBAC (admin vs. user permissions)
  - [ ] Test rate limiting
  - [ ] Test CORS
  - [ ] Test error handling
- [ ] Use `github.com/gin-gonic/gin/httptest` or `net/http/httptest`
- [ ] Mock storage layer for isolation
- [ ] Test with real auth (JWT, API key, Basic)

## API Endpoint Summary

```
# Authentication
POST   /api/auth/login              # Login → JWT + refresh
POST   /api/auth/refresh            # Refresh token
POST   /api/auth/logout             # Logout (blacklist)
GET    /api/auth/verify             # Verify token
POST   /api/auth/api-key            # Create API key
GET    /api/auth/api-key            # List API keys
DELETE /api/auth/api-key/{id}       # Revoke API key

# Users (Admin)
GET    /api/users                   # List users
GET    /api/users/{id}              # Get user
POST   /api/users                   # Create user
PUT    /api/users/{id}              # Update user
DELETE /api/users/{id}              # Delete user

# Calendars
GET    /api/calendars               # List calendars
GET    /api/calendars/{uid}         # Get calendar
POST   /api/calendars               # Create calendar
PUT    /api/calendars/{uid}         # Update calendar
DELETE /api/calendars/{uid}         # Delete calendar
GET    /api/calendars/{uid}/shares  # List shares
POST   /api/calendars/{uid}/shares  # Grant share
DELETE /api/calendars/{uid}/shares/{id}  # Revoke share
GET    /api/calendars/{uid}/export.ics  # Export calendar
GET    /api/calendars/{uid}/stats   # Get stats

# Events
GET    /api/events                  # List events (query params)
GET    /api/events/{uid}            # Get event
POST   /api/events                  # Create event
PUT    /api/events/{uid}            # Update event
DELETE /api/events/{uid}            # Delete event
GET    /api/events/{uid}/export.ics # Export event

# Todos
GET    /api/todos                   # List todos
GET    /api/todos/{uid}             # Get todo
POST   /api/todos                   # Create todo
PUT    /api/todos/{uid}             # Update todo
DELETE /api/todos/{uid}             # Delete todo

# Journals
GET    /api/journals                # List journals
GET    /api/journals/{uid}          # Get journal
POST   /api/journals                # Create journal
PUT    /api/journals/{uid}          # Update journal
DELETE /api/journals/{uid}          # Delete journal

# Free/Busy
GET    /api/freebusy                # Aggregate free/busy
GET    /api/freebusy/{calendar_uid} # Calendar free/busy

# Admin
GET    /api/admin/statistics        # System stats
GET    /api/admin/audit             # Audit log
GET    /api/admin/rate-limits       # Rate limit rules
POST   /api/admin/rate-limits       # Create rate limit
DELETE /api/admin/rate-limits/{id}  # Delete rate limit

# Health
GET    /api/health                  # Health check
GET    /api/status                  # Detailed status
```

## Success Criteria
- [ ] All API endpoints implemented
- [ ] Request validation works
- [ ] Error handling returns consistent responses
- [ ] Auth middleware protects endpoints
- [ ] RBAC middleware enforces permissions
- [ ] Rate limiting works
- [ ] CORS works correctly
- [ ] Pagination works for list endpoints
- [ ] Export endpoints generate valid ICS
- [ ] 85%+ test coverage
- [ ] All tests pass with `go test -race ./app/api/...`

## Notes
- API should be versioned (e.g., `/api/v1/*`) if future changes expected
- Use consistent pagination (page, per_page, total)
- Error codes should be machine-readable (e.g., "calendar.not_found")
- Document API with OpenAPI/Swagger (optional, future phase)

## Next Phase
Move to **Phase 5: CLI Commands** when all items above are complete.
