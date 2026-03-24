# Phase 0: Project Setup and Configuration

## Overview
This phase establishes the project foundation including directory structure, dependencies, configuration management, and development infrastructure.

## Goals
- Set up project directory structure for CalDAV implementation
- Configure Go module and dependencies
- Implement configuration management with Viper
- Set up logging infrastructure
- Create development and production configuration templates
- Establish testing infrastructure
- Create CI/CD configuration
- Set up git workflow with caldav branch

## Tasks

### 1. Directory Structure Setup
- [ ] Create `app/icalendar/` - iCalendar parser/generator
- [ ] Create `app/icalendar/component/` - Component types (VEVENT, VTODO, etc.)
- [ ] Create `app/storage/` - LevelDB storage layer
- [ ] Create `app/storage/calendars/` - Calendar storage
- [ ] Create `app/storage/events/` - Event storage
- [ ] Create `app/storage/todos/` - Todo storage
- [ ] Create `app/storage/journals/` - Journal storage
- [ ] Create `app/storage/freebusy/` - Free/busy generation
- [ ] Create `app/storage/timezones/` - Timezone handling
- [ ] Create `app/storage/recurrence/` - RRULE expansion
- [ ] Create `app/storage/audit/` - Audit logging
- [ ] Create `app/auth/` - Authentication (JWT, API key, Basic)
- [ ] Create `app/auth/jwt/` - JWT management
- [ ] Create `app/auth/apikey/` - API key management
- [ ] Create `app/auth/basic/` - HTTP Basic Auth
- [ ] Create `app/auth/rbac/` - Role-based access control
- [ ] Create `app/api/` - REST API handlers
- [ ] Create `app/api/middleware/` - Auth, RBAC, CORS middleware
- [ ] Create `app/api/handlers/` - API endpoint handlers
- [ ] Create `app/api/models/` - Request/Response models
- [ ] Create `app/webdav/` - WebDAV server integration
- [ ] Create `app/webdav/calendar_fs/` - Calendar filesystem operations
- [ ] Create `app/admin/` - Admin management
- [ ] Create `test/fixtures/` - Test data and fixtures
- [ ] Create `test/fixtures/ics/` - Sample .ics files
- [ ] Create `test/fixtures/config/` - Test configurations

### 2. Dependency Management
- [ ] Review and update `go.mod` with required dependencies
- [ ] Add `github.com/google/leveldb` or `github.com/syndtr/goleveldb`
- [ ] Add `github.com/gorilla/mux` for REST API
- [ ] Add `github.com/golang-jwt/jwt/v5` for JWT
- [ ] Add `github.com/rs/cors` for CORS handling
- [ ] Add `github.com/stretchr/testify` for testing
- [ ] Add `golang.org/x/crypto` for password hashing (already present)
- [ ] Add `github.com/youmark/pkcs8` for private key parsing
- [ ] Update existing dependencies as needed
- [ ] Run `go mod tidy` to clean up

### 3. Configuration Management (Viper)
- [ ] Create `app/config/config.go` - Main configuration struct
- [ ] Implement config loading from YAML file
- [ ] Implement config loading from environment variables
- [ ] Implement config hot reload capability
- [ ] Define configuration sections:
  - [ ] Database (LevelDB paths)
  - [ ] WebDAV (address, port, prefix, TLS)
  - [ ] API (address, port, CORS, auth methods)
  - [ ] Auth (JWT secret, expiry, refresh expiry)
  - [ ] Recurrence (cache expiry, max instances)
  - [ ] Rate Limiting (default limits, admin overrides)
  - [ ] Logging (level, production mode)
  - [ ] Admin (statistics settings, audit settings)
- [ ] Create `config.yaml.example` - Example configuration
- [ ] Create `config.yaml.test` - Test-specific configuration

### 4. Logging Infrastructure
- [ ] Review existing `logrus` setup in `cmd/david/cli/server.go`
- [ ] Create `app/log/logger.go` - Centralized logger factory
- [ ] Implement structured logging for audit trail
- [ ] Implement log rotation configuration
- [ ] Add log levels: debug, info, warn, error, fatal
- [ ] Create log format for production (JSON) vs. development (text)

### 5. Testing Infrastructure
- [ ] Create `test/helpers.go` - Test helpers and utilities
- [ ] Create `test/fixtures/ics/sample_event.ics` - Sample event
- [ ] Create `test/fixtures/ics/sample_todo.ics` - Sample todo
- [ ] Create `test/fixtures/ics/sample_journal.ics` - Sample journal
- [ ] Create `test/fixtures/ics/recurring_event.ics` - Recurring event with RRULE
- [ ] Create `test/fixtures/ics/timezone_event.ics` - Event with VTIMEZONE
- [ ] Create `test/fixtures/ics/invalid.ics` - Invalid/malformed ICS
- [ ] Create `test/coverage.sh` - Coverage reporting script
- [ ] Configure `go test` with race detection
- [ ] Set up table-driven test patterns
- [ ] Set up golden file test patterns
- [ ] Set up fuzz testing infrastructure

### 6. CI/CD Configuration
- [ ] Create `.github/workflows/test.yml` - Automated testing
- [ ] Create `.github/workflows/lint.yml` - Linting checks
- [ ] Create `.github/workflows/build.yml` - Build verification
- [ ] Configure Go version matrix (latest 2 versions)
- [ ] Add coverage reporting to CI
- [ ] Add dependency vulnerability scanning

### 7. Development Tools
- [ ] Review `magefile.go` for build tasks
- [ ] Add `Build` target if not present
- [ ] Add `Test` target for running tests
- [ ] Add `Coverage` target for coverage reports
- [ ] Add `Setup` target for project initialization
- [ ] Create `.gitignore` if not present
- [ ] Add `.editorconfig` for consistent editing
- [ ] Create `Makefile` as alternative to Mage (optional)

### 8. Git Workflow
- [ ] Create `caldav` branch from `main`
- [ ] Document branch workflow in `CONTRIBUTING.md` (optional)
- [ ] Set up commit message conventions
- [ ] Configure pre-commit hooks (optional, for linting)

### 9. Documentation
- [ ] Update `Readme.md` with CalDAV features
- [ ] Create `docs/ARCHITECTURE.md` - System architecture overview
- [ ] Create `docs/API.md` - API endpoint documentation
- [ ] Create `docs/CONFIGURATION.md` - Configuration reference
- [ ] Create `docs/TESTING.md` - Testing guidelines
- [ ] Create `docs/DEPLOYMENT.md` - Deployment guide

## Configuration Structure

### config.yaml Example
```yaml
# Database Configuration
database:
  base_dir: "/var/lib/david"
  calendars_db: "calendars.db"
  events_db: "events.db"
  todos_db: "todos.db"
  journals_db: "journals.db"
  freebusy_db: "freebusy.db"
  timezones_db: "timezones.db"
  users_db: "users.db"
  jwt_blacklist_db: "jwt_blacklist.db"
  audit_db: "audit.db"

# WebDAV Configuration
webdav:
  address: "127.0.0.1"
  port: "8000"
  prefix: "/webdav"
  tls:
    enabled: false
    cert_file: "/etc/david/cert.pem"
    key_file: "/etc/david/key.pem"

# API Configuration
api:
  address: "127.0.0.1"
  port: "8080"
  cors_enabled: true
  cors_origins:
    - "http://localhost:3000"
  auth_methods:
    - "jwt"
    - "api_key"
    - "basic"
  rate_limiting:
    enabled: true
    default_limit: 100
    window: "1m"
    per_ip: true

# Authentication Configuration
auth:
  jwt:
    secret: ""  # Must be set via environment variable
    expiry: "24h"
    refresh_expiry: "7d"
  api_key:
    prefix: "david_"
    expiry: "365d"

# Recurrence Configuration
recurrence:
  cache_expiry_days: 30
  max_cache_instances: 1000

# Logging Configuration
log:
  level: "info"
  production: false
  rotation:
    max_size: 100
    max_backups: 5
    max_age: 30

# Admin Configuration
admin:
  statistics:
    enabled: true
    cache_duration: "5m"
  audit:
    enabled: true
    retention_days: 90
```

### Environment Variables
```bash
# Database
DAVID_DATABASE_BASE_DIR=/var/lib/david

# WebDAV
DAVID_WEBDAV_ADDRESS=127.0.0.1
DAVID_WEBDAV_PORT=8000
DAVID_WEBDAV_PREFIX=/webdav

# API
DAVID_API_ADDRESS=127.0.0.1
DAVID_API_PORT=8080

# Authentication
DAUTH_JWT_SECRET=your-secret-key-here

# Logging
DAVID_LOG_LEVEL=info
DAVID_LOG_PRODUCTION=false
```

## Success Criteria
- [x] Directory structure created
- [x] Dependencies installed and `go mod tidy` passes
- [x] Configuration system working (file + env vars)
- [x] Logging infrastructure in place
- [x] Testing infrastructure ready
- [x] CI/CD pipeline configured
- [x] Development tools configured (Mage, Makefile)
- [x] Git branch `caldav` created
- [x] Documentation templates created
- [x] All tests pass (`go test ./...`)
- [x] Initial coverage baseline established

## Notes
- Phase 0 is foundational - ensure all infrastructure works before proceeding
- Keep configuration flexible for different deployment scenarios
- Test infrastructure should be ready for incremental test writing in later phases
- Document any deviations from the plan in this file

## Next Phase
Move to **Phase 1: Core Parser & Generator** when all items above are complete.
