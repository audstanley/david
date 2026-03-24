# CHANGELOG

All notable changes to david CalDAV Server will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- IANA timezone support with A/B/C download options
- Timezone resolver with priority chain (custom → IANA → UTC)
- Admin notification system with plugin interface
- Dashboard aggregation manager
- Production configuration example
- Development configuration example
- Systemd service file

## [1.0.0] - CalDAV Release - 2026-03-24

### Added
- **Phase 1: iCalendar Parser**
  - Full RFC 2445 compliance
  - Parse VEVENT, VTODO, VJOURNAL, VFREEBUSY, VTIMEZONE
  - Handle recurring events with RRule
  - Timezone support with VTIMEZONE
  - Alarm support (VALARM)
  - 100% test coverage

- **Phase 2: LevelDB Storage**
  - 9 databases: users, calendars, events, todos, journals, freebusy, timezones, recurrence, audit
  - CRUD operations for all entities
  - Indexing for efficient queries
  - Batch operations
  - Backup and restore

- **Phase 3: Authentication**
  - JWT token generation and verification
  - API key generation and management
  - Basic Auth support
  - RBAC with roles: admin, manager, user
  - Permission system
  - Rate limiting
  - Blacklist for token revocation

- **Phase 4: REST API**
  - Gorilla Mux router
  - Auth endpoints (login, refresh, logout)
  - Calendar CRUD endpoints
  - Event CRUD endpoints
  - User management endpoints
  - Error handling
  - CORS support
  - API models (DTOs)

- **Phase 5: CLI Commands**
  - Cobra-based CLI
  - Import/Export commands
  - Calendar management (list, show, create, delete, share)
  - Event management (list, show, delete)
  - User management (list, create, delete)
  - Admin commands (stats, audit)
  - Setup command
  - Server command

- **Phase 6: WebDAV/CalDAV**
  - Custom WebDAV handler
  - PROPFIND for calendar properties
  - REPORT for calendar queries
  - MKCALENDAR method
  - GET/PUT for ICS files
  - DELETE for calendar items
  - Authentication integration

- **Phase 7: Admin & Timezones**
  - System statistics dashboard
  - Notification system with plugin interface
  - IANA tzdata download (3 options)
  - Timezone resolver
  - Systemd service file
  - Configuration examples

### Technical
- Build system with mage
- Comprehensive documentation
- Test infrastructure with fixtures
- Logging with logrus
- Configuration with Viper

### Known Issues
- Some CalDAV REPORT methods are skeleton implementations
- IANA timezone download uses placeholder API
- Notification system requires plugin implementation

[Unreleased]: https://github.com/audstanley/david/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/audstanley/david/releases/tag/v1.0.0
