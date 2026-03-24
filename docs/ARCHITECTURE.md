# david CalDAV Architecture

## Overview

david is a simple WebDAV server extended with full CalDAV (RFC 4791) and iCalendar (RFC 2445) support. This document describes the architecture and design decisions.

## System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    david CLI (Cobra)                    │
├─────────────────────────────────────────────────────────┤
│  david server            (WebDAV + API server)          │
│  david import            (Import ICS → LevelDB)         │
│  david export            (Export LevelDB → ICS)         │
│  david user              (User management)              │
│  david calendar          (Calendar management)          │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│              Multi-Auth REST API (Gorilla Mux)          │
├─────────────────────────────────────────────────────────┤
│  Auth: JWT + API Key + HTTP Basic Auth                 │
│  Roles: Admin, User, Reader, Public                    │
│  Endpoints:                                            │
│    - /api/auth/*         (Login, refresh, verify)       │
│    - /api/users/*        (User CRUD - admin only)       │
│    - /api/calendars/*    (Calendar CRUD + sharing)      │
│    - /api/events/*       (Event CRUD + queries)         │
│    - /api/freebusy/*     (Free/busy generation)         │
│    - /api/export/*       (ICS export endpoints)         │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                  WebDAV Server                           │
├─────────────────────────────────────────────────────────┤
│  /webdav/*               (WebDAV endpoints)             │
│    - PROPFIND, GET, PUT, DELETE, MKCOL                  │
│    - MKCALENDAR (RFC 4791)                              │
│    - REPORT (RFC 4791) - query by date range            │
│    - Sync tokens                                        │
│  Content-Type: text/calendar for .ics files            │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│              LevelDB Storage Layer                       │
├─────────────────────────────────────────────────────────┤
│  users.db              (User accounts + roles)          │
│  calendars.db          (Calendar metadata + shares)     │
│  events.db             (VEVENT + indexes)               │
│  todos.db              (VTODO + indexes)                │
│  journals.db           (VJOURNAL + indexes)             │
│  freebusy.db           (Calculated free/busy)           │
│  timezones.db          (VTIMEZONE definitions)          │
│  jwt_blacklist.db      (Revoked tokens)                 │
│  audit.db              (Audit trail)                    │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                  Frontend (Optional)                     │
├─────────────────────────────────────────────────────────┤
│  Svelte/React/Vue app consuming /api/* endpoints        │
│  (Built separately, not part of this project)           │
└─────────────────────────────────────────────────────────┘
```

## Component Descriptions

### iCalendar Parser/Generator (`app/icalendar/`)

The core parser implements RFC 2445 compliance for parsing and generating iCalendar files.

**Key Features:**
- Full RFC 2445 content line parsing
- Line folding/unfolding support
- All value types (DATE, DATE-TIME, DURATION, RRULE, etc.)
- Property parameter parsing
- Component parsing (VEVENT, VTODO, VJOURNAL, VFREEBUSY, VTIMEZONE, VALARM)
- Bidirectional serialization (ICS ↔ Go structs)

### Storage Layer (`app/storage/`)

Multi-Database LevelDB setup for efficient data storage and retrieval.

**Database Structure:**
- **users.db** - User accounts, roles, API keys
- **calendars.db** - Calendar metadata, shares, access control
- **events.db** - Event data with date-indexed keys
- **todos.db** - Todo items
- **journals.db** - Journal entries
- **freebusy.db** - Calculated free/busy time blocks
- **timezones.db** - IANA tzdata and custom VTIMEZONE
- **jwt_blacklist.db** - Revoked tokens
- **audit.db** - Audit trail

**Index Strategy:**
- UID lookup: O(1) direct key lookup
- Date range queries: Range scan on lexicographically sortable date keys
- Secondary indexes: Optional (attendee, status, priority)

### Authentication (`app/auth/`)

Multi-layer authentication supporting JWT, API keys, and HTTP Basic Auth.

**Features:**
- JWT with refresh tokens
- API key generation and verification
- HTTP Basic Auth integration
- Role-based access control (Admin, User, Reader)
- Token blacklist for revocation
- Rate limiting (per IP, per user, admin overrides)

### REST API (`app/api/`)

Gorilla Mux-based REST API for programmatic access.

**Endpoints:**
- `/api/auth/*` - Authentication operations
- `/api/users/*` - User management (admin)
- `/api/calendars/*` - Calendar CRUD and sharing
- `/api/events/*` - Event CRUD and queries
- `/api/todos/*` - Todo management
- `/api/journals/*` - Journal management
- `/api/freebusy/*` - Free/busy generation
- `/api/export/*` - ICS export
- `/api/admin/*` - Admin operations

### WebDAV/CalDAV (`app/webdav/`)

RFC 4791 CalDAV extension for standard calendar client support.

**Supported Methods:**
- PROPFIND - Property discovery
- MKCALENDAR - Create calendar (RFC 4791)
- REPORT - Query calendar (RFC 4791)
- PROPPATCH - Modify properties
- GET/PUT/DELETE - ICS file operations

**Calendar Properties:**
- C:calendar-collection
- C:calendar-home-set
- C:calendar-description
- C:calendar-color
- C:supported-calendar-component-set
- Sync tokens

### CLI (`cmd/david/cli/`)

Cobra-based CLI for calendar management.

**Commands:**
- `david server` - Start WebDAV + API server
- `david import` - Import ICS files
- `david export` - Export to ICS
- `david calendar` - Calendar management
- `david event` - Event management
- `david user` - User management
- `david admin` - Admin operations

## Data Flow

### ICS Import Flow
1. User uploads ICS file via WebDAV or CLI
2. Parser validates and parses ICS
3. Components stored in appropriate LevelDB tables
4. Indexes updated for efficient queries
5. Recurrence cache populated if applicable

### Event Query Flow
1. Client queries by date range
2. Query translates to LevelDB range scan
3. Results expanded for recurrence instances (from cache)
4. Results returned in ICS format

### Auth Flow
1. Client sends credentials (Basic, API key, or JWT)
2. Auth manager validates credentials
3. Token generated (if JWT)
4. Permissions checked via RBAC
5. Request proceeds with user context

## Design Decisions

### Why LevelDB?
- Embedded, no external dependencies
- Fast key-value operations
- Good for high-throughput writes
- Simple API

### Why Hybrid Recurrence?
- Store RRULE once (efficient storage)
- Pre-compute cache for near-term (fast queries)
- On-demand expansion for distant future
- Cache invalidation on rule changes

### Why Multiple Auth Methods?
- JWT for web/mobile clients
- API keys for programmatic access
- Basic auth for simplicity/compatibility

### Why CalDAV + REST?
- CalDAV for standard client support
- REST for custom frontend development
- Both can coexist and use same backend

## Testing Strategy

- Table-driven unit tests
- Golden file tests for ICS parsing
- Fuzz testing for parser robustness
- Integration tests with real LevelDB
- 85%+ coverage target

## Performance Considerations

- UID lookups: O(1) via direct keys
- Date range queries: Efficient via sorted keys
- Recurrence: Cached for common queries
- Rate limiting: Prevents abuse
- Audit logging: Async where possible
