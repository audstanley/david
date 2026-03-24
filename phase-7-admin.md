# Phase 7: Admin, IANA Timezone & Final Integration

## Overview
This phase implements admin management features, IANA timezone integration, and brings all phases together for final integration and deployment.

## Goals
- Implement admin management CLI and API
- Implement IANA tzdata download and usage
- Support custom VTIMEZONE overrides
- Create systemd/init scripts
- Create configuration examples
- Final integration testing
- Documentation completion
- Achieve overall 85%+ coverage

## Tasks

### 1. Admin Management Enhancements
- [ ] Review `app/admin/` from Phase 5
- [ ] Create `app/admin/dashboard.go` - Dashboard data aggregation
  - [ ] Aggregate statistics across all data
  - [ ] Cache results (configurable duration)
  - [ ] Support real-time updates (optional)
- [ ] Create `app/admin/notifications.go` - Admin notifications
  - [ ] Send alerts for system issues
  - [ ] Storage capacity warnings
  - [ ] Failed login attempts
  - [ ] Audit anomalies
- [ ] Implement admin API endpoints (Phase 4)
  - [ ] GET /api/admin/statistics
  - [ ] GET /api/admin/audit
  - [ ] GET /api/admin/rate-limits
  - [ ] POST /api/admin/rate-limits
  - [ ] DELETE /api/admin/rate-limits/{id}

### 2. IANA Timezone Integration
- [ ] Create `app/storage/timezones/iana.go` - IANA tzdata
  - [ ] Download IANA tzdata (from timezonedb.com API or similar)
  - [ ] Parse IANA timezone database format
  - [ ] Extract timezone rules (offsets, DST)
  - [ ] Cache downloaded data in LevelDB
  - [ ] Implement automatic updates (periodic check)
- [ ] Create `app/storage/timezones/iana_cache.go` - Cache management
  - [ ] Store cached tzdata in LevelDB
  - [ ] Check for updates (daily/weekly)
  - [ ] Handle download failures (fallback to cached)
  - [ ] Implement cache invalidation
- [ ] Create `app/storage/timezones/resolver.go` - Timezone resolution
  - [ ] Resolve TZID to timezone rules
  - [ ] First check custom VTIMEZONE (Phase 2)
  - [ ] Then check cached IANA tzdata
  - [ ] Fall back to UTC if not found
  - [ ] Apply timezone rules to datetime
  - [ ] Convert local time ↔ UTC
  - [ ] Handle DST transitions
- [ ] Create `app/storage/timezones/types.go` - Timezone types
  - [ ] TimezoneRule struct (offsets, DST rules)
  - [ ] TimezoneOffset struct (from, to, transition)

### 3. Custom VTIMEZONE Management
- [ ] Review `app/storage/timezones/store.go` (Phase 2)
- [ ] Enhance VTIMEZONE CRUD operations
- [ ] Implement VTIMEZONE upload from ICS
- [ ] Implement VTIMEZONE validation
- [ ] Implement VTIMEZONE conflict resolution (custom vs. IANA)
- [ ] Create `app/api/handlers/timezones.go` - Timezone API
  - [ ] GET /api/timezones - List available timezones
  - [ ] GET /api/timezones/{tzid} - Get timezone definition
  - [ ] POST /api/timezones - Upload custom VTIMEZONE
  - [ ] DELETE /api/timezones/{tzid} - Delete custom VTIMEZONE

### 4. Systemd/Init Scripts
- [ ] Create `contrib/systemd/david.service` - Systemd service
  - [ ] Define service configuration
  - [ ] Set restart policy
  - [ ] Configure logging (journal)
  - [ ] Set resource limits
  - [ ] Define dependencies (network, filesystem)
- [ ] Create `contrib/init/david.init` - SysV init script (optional)
  - [ ] Traditional init script format
  - [ ] Start/stop/restart functions
  - [ ] Status command
- [ ] Create `contrib/environment/david.env` - Environment file
  - [ ] Document all environment variables
  - [ ] Default values
  - [ ] Examples

### 5. Configuration Examples
- [ ] Update `config.yaml.example` - Comprehensive example
  - [ ] All configuration sections
  - [ ] Comments explaining each option
  - [ ] Example values
- [ ] Create `config.yaml.production` - Production config example
  - [ ] Optimized settings
  - [ ] Security best practices
  - [ ] TLS configuration
- [ ] Create `config.yaml.development` - Development config example
  - [ ] Debug settings
  - [ ] Development-specific options

### 6. Documentation
- [ ] Update `Readme.md`
  - [ ] Add CalDAV support section
  - [ ] Add API documentation link
  - [ ] Add admin features section
  - [ ] Add deployment section
  - [ ] Update installation instructions
- [ ] Create `docs/CALDAV.md` - CalDAV guide
  - [ ] CalDAV overview
  - [ ] Client configuration examples
  - [ ] Troubleshooting
- [ ] Create `docs/API.md` - API documentation
  - [ ] All API endpoints
  - [ ] Request/response examples
  - [ ] Authentication guide
  - [ ] Error codes
- [ ] Create `docs/ADMIN.md` - Admin guide
  - [ ] User management
  - [ ] Calendar management
  - [ ] Statistics and monitoring
  - [ ] Audit logging
- [ ] Create `docs/TIMEZONES.md` - Timezone guide
  - [ ] IANA tzdata usage
  - [ ] Custom VTIMEZONE
  - [ ] Timezone conversion
- [ ] Create `docs/DEPLOYMENT.md` - Deployment guide
  - [ ] Binary deployment
  - [ ] systemd setup
  - [ ] TLS configuration
  - [ ] Reverse proxy setup
- [ ] Create `docs/SECURITY.md` - Security guide
  - [ ] Authentication options
  - [ ] Rate limiting
  - [ ] Audit logging
  - [ ] Best practices

### 7. Final Integration Testing
- [ ] Create `test/integration/full_test.go` - Full workflow tests
  - [ ] Test complete user journey
  - [ ] Test admin workflow
  - [ ] Test CalDAV client workflow
  - [ ] Test import/export workflow
  - [ ] Test public calendar workflow
- [ ] Test with real CalDAV clients
  - [ ] Apple Calendar (if available on test system)
  - [ ] Thunderbird Lightning
  - [ ] DAVx5 (Android)
  - [ ] vdirsyncer (command-line)
- [ ] Performance testing
  - [ ] Load test with many events
  - [ ] Test recurrence expansion performance
  - [ ] Test date range query performance
- [ ] Coverage analysis
  - [ ] Run `go test -cover`
  - [ ] Verify 85%+ overall coverage
  - [ ] Identify gaps
  - [ ] Add tests for gaps

### 8. Build & Release
- [ ] Update `magefile.go`
  - [ ] Add `Build` target if not present
  - [ ] Add `Release` target (build binaries)
  - [ ] Add `TestAll` target (all tests + coverage)
  - [ ] Add `Lint` target
- [ ] Create `scripts/build.sh` - Build script
  - [ ] Build binary with version info
  - [ ] Create release artifacts
  - [ ] Generate checksums
- [ ] Create `CHANGELOG.md` - Change log
  - [ ] Track all changes across phases
  - [ ] Format: Keep a Changelog style
- [ ] Create `VERSION` file - Version tracking
- [ ] Tag release in git (v1.0.0-caldav)

## Data Models

### Admin Statistics
```go
type SystemStats struct {
    TotalUsers      int64
    TotalCalendars  int64
    TotalEvents     int64
    TotalTodos      int64
    TotalJournals   int64
    TotalStorage    int64  // bytes
    Uptime          time.Duration
    LastBackup      time.Time
}

type UserStats struct {
    UserID         string
    Username       string
    Role           string
    CalendarsOwned int64
    EventsCreated  int64
    LastLogin      time.Time
}

type CalendarStats struct {
    CalendarUID  string
    DisplayName  string
    EventCount   int64
    TodoCount    int64
    JournalCount int64
    StorageSize  int64
    LastModified time.Time
}
```

## Success Criteria
- [ ] Admin features work (stats, audit, ratelimit)
- [ ] IANA tzdata downloads and caches
- [ ] Custom VTIMEZONE works alongside IANA
- [ ] Timezone resolution handles DST correctly
- [ ] Systemd service file works
- [ ] Configuration examples are complete
- [ ] Documentation is comprehensive
- [ ] Full integration tests pass
- [ ] CalDAV clients can connect
- [ ] 85%+ overall coverage
- [ ] Build produces working binary

## Deployment Checklist
- [ ] Binary built and tested
- [ ] Systemd service file in place
- [ ] Configuration file created
- [ ] TLS certificates configured (if using TLS)
- [ ] Reverse proxy configured (if using)
- [ ] Firewall rules configured
- [ ] Backup strategy in place
- [ ] Monitoring configured
- [ ] Logging configured
- [ ] Admin user created
- [ ] Documentation accessible

## Release Notes Template

```markdown
# Version X.Y.Z - CalDAV Release

## New Features
- Full CalDAV support (RFC 4791)
- REST API with JWT, API Key, and Basic Auth
- Admin management interface
- IANA timezone support
- Calendar sharing with role-based access
- Rate limiting
- Audit logging

## Improvements
- Enhanced iCalendar parsing
- Optimized date range queries
- Recurrence cache for better performance

## Breaking Changes
- None (backward compatible)

## Known Issues
- [List any known issues]

## Migration Guide
- [If applicable]
```

## Next Phase
Move to **Phase 8: Testing & Polish** if additional testing is needed, or proceed to release.
