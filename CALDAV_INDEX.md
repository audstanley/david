# david CalDAV Implementation - Master Index

## Overview
This repository contains the implementation plan for adding full CalDAV (RFC 4791) and iCalendar (RFC 2445) support to the david WebDAV server.

## Branch
**caldav** - All work is done on this branch

## Phase Files

| Phase | File | Status | Description |
|-------|------|--------|-------------|
| 0 | [phase-0-setup.md](phase-0-setup.md) | ⏳ Pending | Project setup and configuration |
| 1 | [phase-1-parser.md](phase-1-parser.md) | ⏳ Pending | iCalendar parser/generator |
| 2 | [phase-2-storage.md](phase-2-storage.md) | ⏳ Pending | LevelDB storage layer |
| 3 | [phase-3-auth.md](phase-3-auth.md) | ⏳ Pending | Authentication & authorization |
| 4 | [phase-4-api.md](phase-4-api.md) | ⏳ Pending | REST API (Gorilla Mux) |
| 5 | [phase-5-cli.md](phase-5-cli.md) | ⏳ Pending | CLI commands (Cobra) |
| 6 | [phase-6-webdav.md](phase-6-webdav.md) | ⏳ Pending | WebDAV/CalDAV integration |
| 7 | [phase-7-admin.md](phase-7-admin.md) | ⏳ Pending | Admin features & deployment |

## Project Structure (Target)

```
david/
├── phase-0-setup.md
├── phase-1-parser.md
├── phase-2-storage.md
├── phase-3-auth.md
├── phase-4-api.md
├── phase-5-cli.md
├── phase-6-webdav.md
├── phase-7-admin.md
├── app/
│   ├── app.go
│   ├── config/
│   │   └── config.go
│   ├── log/
│   │   └── logger.go
│   ├── icalendar/
│   │   ├── parser.go
│   │   ├── generator.go
│   │   ├── types.go
│   │   ├── properties.go
│   │   ├── value_types.go
│   │   ├── parameter.go
│   │   ├── validation.go
│   │   ├── errors.go
│   │   ├── rrule.go
│   │   └── component/
│   │       ├── vcalendar.go
│   │       ├── vevent.go
│   │       ├── vtodo.go
│   │       ├── vjournal.go
│   │       ├── vfreebusy.go
│   │       ├── vtimezone.go
│   │       ├── valarm.go
│   │       └── common.go
│   ├── storage/
│   │   ├── leveldb.go
│   │   ├── interface.go
│   │   ├── errors.go
│   │   ├── users/
│   │   │   ├── store.go
│   │   │   └── auth.go
│   │   ├── calendars/
│   │   │   ├── store.go
│   │   │   ├── indexes.go
│   │   │   └── shares.go
│   │   ├── events/
│   │   │   ├── store.go
│   │   │   ├── indexes.go
│   │   │   └── query.go
│   │   ├── todos/
│   │   │   ├── store.go
│   │   │   └── query.go
│   │   ├── journals/
│   │   │   ├── store.go
│   │   │   └── query.go
│   │   ├── freebusy/
│   │   │   ├── generator.go
│   │   │   └── query.go
│   │   ├── timezones/
│   │   │   ├── store.go
│   │   │   ├── iana.go
│   │   │   ├── resolver.go
│   │   │   └── cache.go
│   │   ├── recurrence/
│   │   │   ├── cache.go
│   │   │   ├── expand.go
│   │   │   └── validator.go
│   │   └── audit/
│   │       ├── store.go
│   │       ├── types.go
│   │       └── query.go
│   ├── auth/
│   │   ├── manager.go
│   │   ├── types.go
│   │   ├── jwt/
│   │   │   ├── token.go
│   │   │   ├── refresh.go
│   │   │   └── blacklist.go
│   │   ├── apikey/
│   │   │   ├── generator.go
│   │   │   ├── verifier.go
│   │   │   └── store.go
│   │   ├── basic/
│   │   │   ├── handler.go
│   │   │   └── verifier.go
│   │   └── rbac/
│   │       ├── roles.go
│   │       ├── permissions.go
│   │       └── checker.go
│   ├── api/
│   │   ├── router.go
│   │   ├── errors.go
│   │   ├── models/
│   │   │   ├── response.go
│   │   │   ├── request.go
│   │   │   └── validation.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── rbac.go
│   │   │   ├── ratelimit.go
│   │   │   ├── cors.go
│   │   │   ├── logging.go
│   │   │   └── recovery.go
│   │   └── handlers/
│   │       ├── auth.go
│   │       ├── users.go
│   │       ├── calendars.go
│   │       ├── shares.go
│   │       ├── events.go
│   │       ├── todos.go
│   │       ├── journals.go
│   │       ├── freebusy.go
│   │       ├── export.go
│   │       ├── admin.go
│   │       └── health.go
│   ├── webdav/
│   │   ├── handler.go
│   │   ├── filesystem.go
│   │   ├── resolver.go
│   │   ├── content_type.go
│   │   ├── import_handler.go
│   │   ├── export_handler.go
│   │   ├── auth.go
│   │   └── calendar_fs/
│   │       ├── mkcalendar.go
│   │       ├── propfind.go
│   │       ├── report.go
│   │       ├── sync.go
│   │       └── types.go
│   └── admin/
│       ├── dashboard.go
│       ├── notifications.go
│       └── timezones.go
├── cmd/
│   └── david/
│       ├── main.go
│       └── cli/
│           ├── root.go
│           ├── server.go
│           ├── import.go
│           ├── export.go
│           ├── calendar/
│           │   ├── list.go
│           │   ├── show.go
│           │   ├── create.go
│           │   ├── update.go
│           │   ├── delete.go
│           │   └── share.go
│           ├── event/
│           │   ├── list.go
│           │   ├── show.go
│           │   └── delete.go
│           ├── user/
│           │   ├── list.go
│           │   ├── create.go
│           │   ├── delete.go
│           │   ├── role.go
│           │   └── apikey.go
│           └── admin/
│               ├── stats.go
│               ├── audit.go
│               └── ratelimit.go
├── test/
│   ├── fixtures/
│   │   ├── ics/
│   │   │   ├── event_basic.ics
│   │   │   ├── event_recurring.ics
│   │   │   ├── event_timezone.ics
│   │   │   ├── event_alarm.ics
│   │   │   ├── todo_basic.ics
│   │   │   ├── todo_completed.ics
│   │   │   ├── journal_basic.ics
│   │   │   ├── freebusy_basic.ics
│   │   │   ├── timezone_america_new_york.ics
│   │   │   └── malformed.ics
│   │   ├── config/
│   │   │   └── test.yaml
│   │   └── cli/
│   │       └── import/
│   └── helpers.go
├── contrib/
│   ├── systemd/
│   │   └── david.service
│   ├── init/
│   │   └── david.init
│   └── environment/
│       └── david.env
├── docs/
│   ├── ARCHITECTURE.md
│   ├── API.md
│   ├── CONFIGURATION.md
│   ├── TESTING.md
│   ├── DEPLOYMENT.md
│   ├── CALDAV.md
│   ├── ADMIN.md
│   ├── TIMEZONES.md
│   └── SECURITY.md
├── config.yaml.example
├── config.yaml.test
├── README.md (updated)
├── CHANGELOG.md
├── VERSION
├── go.mod
├── go.sum
├── magefile.go
├── .gitignore
└── .editorconfig
```

## Implementation Order

1. ✅ Create phase plan files (current)
2. ⏳ Phase 0: Project Setup
3. ⏳ Phase 1: iCalendar Parser
4. ⏳ Phase 2: LevelDB Storage
5. ⏳ Phase 3: Authentication
6. ⏳ Phase 4: REST API
7. ⏳ Phase 5: CLI Commands
8. ⏳ Phase 6: WebDAV/CalDAV
9. ⏳ Phase 7: Admin & Deployment

## Quick Start

```bash
# Switch to caldav branch
git checkout -b caldav

# Start with Phase 0
# Read phase-0-setup.md for detailed instructions
```

## Testing

```bash
# Run all tests
go test ./... -race

# Run with coverage
go test ./... -race -coverprofile=coverage.out
go tool cover -html=coverage.out

# Check specific phase
go test ./app/icalendar/... -v -cover
```

## Configuration

```bash
# Example configuration
cat config.yaml.example

# Set JWT secret (required)
export DAVID_JWT_SECRET="your-secret-key"

# Run server
david server --config config.yaml
```

## API Access

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Authorization: Basic base64(user:password)"

# List calendars
curl http://localhost:8080/api/calendars \
  -H "Authorization: Bearer <token>"
```

## CalDAV Client Configuration

```
# Apple Calendar / macOS
Server: https://yourserver.com/webdav/
Username: your_username
Password: your_password

# Thunderbird / Lightning
Server: https://yourserver.com/webdav/
Username: your_username
Password: your_password
```

## License
Apache 2.0 (same as parent project)
