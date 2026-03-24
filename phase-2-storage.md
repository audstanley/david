# Phase 2: LevelDB Storage Layer

## Overview
This phase implements the LevelDB storage layer for all iCalendar data. This provides efficient CRUD operations, optimized indexing for queries, and persistence for calendars, events, todos, journals, and related data.

## Goals
- Implement multi-Database LevelDB setup (7 databases)
- Create optimized indexes for UID lookups and date range queries
- Implement CRUD operations for all data types
- Implement recurrence cache with hybrid expansion
- Implement timezone resolution (IANA + custom)
- Implement audit logging
- Achieve 90%+ unit test coverage

## Tasks

### 1. LevelDB Infrastructure
- [ ] Create `app/storage/leveldb.go` - Database factory and management
  - [ ] Create multi-DB manager (open, close, get DB)
  - [ ] Implement connection pooling
  - [ ] Implement database path management
  - [ ] Add database existence checks
  - [ ] Implement backup/restore utilities
- [ ] Create `app/storage/interface.go` - Storage interface abstraction
  - [ ] Define CalendarStore interface
  - [ ] Define EventStore interface
  - [ ] Define TodoStore interface
  - [ ] Define JournalStore interface
  - [ ] Define FreeBusyStore interface
  - [ ] Define TimeZoneStore interface
  - [ ] Define UserStore interface
- [ ] Create `app/storage/errors.go` - Custom storage errors

### 2. User Storage
- [ ] Create `app/storage/users/store.go` - User CRUD operations
  - [ ] Create user (with password hash)
  - [ ] Get user by ID
  - [ ] Get user by username
  - [ ] Update user
  - [ ] Delete user
  - [ ] List users (with pagination)
- [ ] Create `app/storage/users/auth.go` - Authentication
  - [ ] Verify password hash (bcrypt, argon2, scrypt)
  - [ ] Hash new password
  - [ ] Get user by auth token
- [ ] Create `app/storage/users/types.go` - User data model
  - [ ] UserID, Username, PasswordHash
  - [ ] Role (admin, user)
  - [ ] Email, DisplayName
  - [ ] Created, Updated timestamps
  - [ ] APIKeys (list)

### 3. Calendar Storage
- [ ] Create `app/storage/calendars/store.go` - Calendar CRUD
  - [ ] Create calendar (UID, owner_id, display_name, description, color, is_public)
  - [ ] Get calendar by UID
  - [ ] Get calendar by name
  - [ ] Update calendar
  - [ ] Delete calendar
  - [ ] List calendars (by owner, with filters)
- [ ] Create `app/storage/calendars/indexes.go` - Index management
  - [ ] Index: cal:uid:{uid} → Calendar metadata
  - [ ] Index: cal:owner:{owner_id} → List of calendar UIDs
  - [ ] Index: cal:name:{name} → UID lookup
  - [ ] Index: cal:public:true → List of public calendars
- [ ] Create `app/storage/calendars/shares.go` - Share management
  - [ ] Create share (calendar_uid, target_user_id, role: read/write/admin/everyone)
  - [ ] Get shares for calendar
  - [ ] Delete share
  - [ ] Check access (read/write/admin) for user
  - [ ] List all users with access
- [ ] Create `app/storage/calendars/types.go` - Calendar data model

### 4. Event Storage
- [ ] Create `app/storage/events/store.go` - Event CRUD
  - [ ] Create event (full event struct)
  - [ ] Get event by UID
  - [ ] Update event (with sequence number)
  - [ ] Delete event
  - [ ] List events (by calendar, with filters)
- [ ] Create `app/storage/events/indexes.go` - Index management
  - [ ] Primary: comp:uid:{uid} → Event (JSON)
  - [ ] Date index: comp:date:{cal_uid}:{YYYYMMDDTHHMMSSZ}:{seq}:{uid} → Event
  - [ ] UID index: comp:calendar:{cal_uid}:uid:{uid} → Event
  - [ ] Attendee index: comp:attendee:{cal_uid}:{email}:{uid} → Event (optional)
  - [ ] Status index: comp:status:{cal_uid}:{status}:{uid} → Event (optional)
- [ ] Create `app/storage/events/query.go` - Query logic
  - [ ] Get event by UID (O(1) lookup)
  - [ ] Get events in date range (range scan on date index)
  - [ ] Get events by calendar + date range
  - [ ] Get recurring instances (from cache)
  - [ ] Get expanded recurrence instances
- [ ] Create `app/storage/events/types.go` - Event data model

### 5. Todo Storage
- [ ] Create `app/storage/todos/store.go` - Todo CRUD
  - [ ] Create todo
  - [ ] Get todo by UID
  - [ ] Update todo
  - [ ] Delete todo
  - [ ] List todos (by calendar, with filters)
- [ ] Create `app/storage/todos/indexes.go` - Index management
  - [ ] Primary: comp:uid:{uid} → Todo
  - [ ] Date index: comp:date:{cal_uid}:{date}:{seq}:{uid} → Todo
- [ ] Create `app/storage/todos/query.go` - Query logic
- [ ] Create `app/storage/todos/types.go` - Todo data model

### 6. Journal Storage
- [ ] Create `app/storage/journals/store.go` - Journal CRUD
  - [ ] Create journal
  - [ ] Get journal by UID
  - [ ] Update journal
  - [ ] Delete journal
  - [ ] List journals
- [ ] Create `app/storage/journals/indexes.go` - Index management
- [ ] Create `app/storage/journals/query.go` - Query logic
- [ ] Create `app/storage/journals/types.go` - Journal data model

### 7. Free/Busy Storage
- [ ] Create `app/storage/freebusy/generator.go` - Free/busy calculation
  - [ ] Generate free/busy for single calendar (from events)
  - [ ] Generate free/busy for multiple calendars
  - [ ] Handle transparency (TRANSPARENT vs OPAQUE)
  - [ ] Handle conflicts
  - [ ] Calculate busy blocks from events
  - [ ] Calculate free blocks (inverted)
- [ ] Create `app/storage/freebusy/query.go` - Query logic
  - [ ] Get free/busy for date range
  - [ ] Aggregate free/busy across multiple calendars
  - [ ] Find available time slots
- [ ] Create `app/storage/freebusy/types.go` - Free/busy data model
- [ ] Create `app/storage/freebusy/store.go` - Persist free/busy (optional)

### 8. Timezone Storage
- [ ] Create `app/storage/timezones/store.go` - Custom VTIMEZONE
  - [ ] Store custom VTIMEZONE definition
  - [ ] Get VTIMEZONE by TZID
  - [ ] Update VTIMEZONE
  - [ ] Delete VTIMEZONE
- [ ] Create `app/storage/timezones/iana.go` - IANA tzdata
  - [ ] Download IANA tzdata (from timezonedb.com or similar)
  - [ ] Parse IANA timezone database
  - [ ] Cache downloaded data
  - [ ] Fallback to UTC if not found
- [ ] Create `app/storage/timezones/resolver.go` - Timezone resolution
  - [ ] Resolve TZID to UTC offset
  - [ ] Handle VTIMEZONE components (STANDARD/DAYLIGHT)
  - [ ] Apply RRULE-based timezone rules
  - [ ] Convert local time to UTC
  - [ ] Convert UTC to local time
- [ ] Create `app/storage/timezones/types.go` - Timezone data model

### 9. Recurrence Cache
- [ ] Create `app/storage/recurrence/cache.go` - Cache management
  - [ ] Store expanded recurrence instances
  - [ ] Cache versioning (increment on RRULE change)
  - [ ] Cache expiry (based on configured days)
  - [ ] Cache cleanup (remove expired entries)
- [ ] Create `app/storage/recurrence/expand.go` - RRULE expansion
  - [ ] Expand RRULE to instances
  - [ ] Handle FREQ=SECONDLY, MINUTELY, HOURLY
  - [ ] Handle FREQ=DAILY, WEEKLY
  - [ ] Handle FREQ=MONTHLY, YEARLY
  - [ ] Apply EXDATE exceptions
  - [ ] Apply RDATE additions
  - [ ] Limit expansion to cache window
- [ ] Create `app/storage/recurrence/validator.go` - Validation
  - [ ] Validate RRULE syntax
  - [ ] Validate RRULE constraints (e.g., BYDAY with FREQ=WEEKLY)
  - [ ] Validate recurrence date ranges
  - [ ] Detect infinite recurrences

### 10. Audit Logging
- [ ] Create `app/storage/audit/store.go` - Audit log operations
  - [ ] Log user actions (create, update, delete)
  - [ ] Log calendar actions (create, update, delete, share)
  - [ ] Log event actions (create, update, delete)
  - [ ] Log admin actions (user management, system config)
- [ ] Create `app/storage/audit/types.go` - Audit entry model
  - [ ] Timestamp, UserID, Action, EntityType, EntityID, Details
- [ ] Create `app/storage/audit/query.go` - Audit queries
  - [ ] Get audit log by user
  - [ ] Get audit log by entity
  - [ ] Get audit log by date range
  - [ ] Get audit log by action type
- [ ] Implement audit log rotation (retention days)

### 11. Unit Tests
- [ ] Create `app/storage/leveldb_test.go`
  - [ ] Test database creation
  - [ ] Test database opening
  - [ ] Test multi-DB management
- [ ] Create `app/storage/users/store_test.go`
  - [ ] Test CRUD operations
  - [ ] Test password hashing
- [ ] Create `app/storage/users/auth_test.go`
  - [ ] Test password verification
  - [ ] Test token-based user lookup
- [ ] Create `app/storage/calendars/store_test.go`
- [ ] Create `app/storage/calendars/shares_test.go`
  - [ ] Test share creation
  - [ ] Test access control (read/write/admin)
- [ ] Create `app/storage/events/store_test.go`
- [ ] Create `app/storage/events/query_test.go`
  - [ ] Test UID lookup
  - [ ] Test date range queries
  - [ ] Test recurrence expansion
- [ ] Create `app/storage/todos/store_test.go`
- [ ] Create `app/storage/journals/store_test.go`
- [ ] Create `app/storage/freebusy/generator_test.go`
- [ ] Create `app/storage/timezones/resolver_test.go`
  - [ ] Test timezone conversion
  - [ ] Test IANA tzdata usage
  - [ ] Test custom VTIMEZONE
- [ ] Create `app/storage/recurrence/expand_test.go`
  - [ ] Test RRULE expansion for all frequencies
  - [ ] Test EXDATE handling
  - [ ] Test RDATE handling
- [ ] Create `app/storage/recurrence/cache_test.go`
- [ ] Create `app/storage/audit/store_test.go`

### 12. Integration Tests
- [ ] Create `app/storage/integration_test.go`
  - [ ] Test full workflow: create calendar → add event → query events
  - [ ] Test share workflow: create calendar → share with user → verify access
  - [ ] Test recurrence workflow: create recurring event → query instances
  - [ ] Test timezone workflow: create event with TZID → verify UTC conversion

## Data Models

### Calendar
```go
type Calendar struct {
    UID         string
    OwnerID     string
    DisplayName string
    Description string
    Color       string
    IsPublic    bool
    PublicHash  string  // Obscured URL hash
    Timezone    string
    Created     time.Time
    Updated     time.Time
    Shares      []Share
}

type Share struct {
    ID          string
    CalendarUID string
    TargetUserID string
    Role        string  // read, write, admin, everyone
    Created     time.Time
    CreatedBy   string
}
```

### Event
```go
type Event struct {
    UID          string
    CalendarUID  string
    DTStamp      time.Time
    DTStart      time.Time  // Stored in UTC
    DTEnd        time.Time  // Stored in UTC
    Duration     *time.Duration
    Summary      string
    Description  string
    Location     string
    Organizer    string
    Attendees    []Attendee
    Status       string  // CONFIRMED, TENTATIVE, CANCELLED
    Class        string  // PUBLIC, PRIVATE, CONFIDENTIAL
    Priority     int
    Sequence     int
    Categories   []string
    Resources    []string
    RRule        *RRule
    ExDates      []time.Time
    Created      time.Time
    LastModified time.Time
    Alarms       []Alarm
}

type Attendee struct {
    Email     string
    Name      string
    RSVP      bool
    Role      string  // REQ-PARTICIPANT, OPT-PARTICIPANT, NON-PARTICIPANT
    PartStat  string  // NEEDS-ACTION, ACCEPTED, DECLINED, TENTATIVE, DELEGATED
    Cutype    string  // INDIVIDUAL, GROUP, RESOURCE, ROOM
}
```

### Audit Entry
```go
type AuditEntry struct {
    ID        string
    Timestamp time.Time
    UserID    string
    Action    string  // create, update, delete, share_grant, share_revoke
    EntityType string // user, calendar, event, todo, journal
    EntityID  string
    Details   string  // JSON-serialized details
}
```

## Success Criteria
- [ ] All 7 databases created and accessible
- [ ] CRUD operations work for all data types
- [ ] UID lookup is O(1) (direct key lookup)
- [ ] Date range queries work efficiently (range scan)
- [ ] Recurrence cache works correctly
- [ ] RRULE expansion handles all frequencies
- [ ] Timezone resolution works (IANA + custom)
- [ ] Share management works (create, read, delete)
- [ ] Access control works (read/write/admin roles)
- [ ] Audit logging captures all actions
- [ ] 90%+ unit test coverage
- [ ] All tests pass with `go test -race ./app/storage/...`

## Notes
- LevelDB keys are bytes - use consistent encoding (e.g., UTF-8 strings)
- Store complex objects as JSON for flexibility
- Index keys should be lexicographically sortable for range queries
- Cache should expire automatically (implement cleanup job)
- Audit log should never fail (buffer writes if needed)

## Next Phase
Move to **Phase 3: Authentication & Authorization** when all items above are complete.
