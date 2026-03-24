# Phase 6: WebDAV Integration

## Overview
This phase integrates CalDAV support into the WebDAV server. It implements RFC 4791 (CalDAV) extensions including MKCALENDAR, REPORT, PROPFIND for calendar properties, and sync tokens. WebDAV provides standard calendar access for clients like Apple Calendar, Thunderbird, and others.

## Goals
- Implement CalDAV extensions (RFC 4791)
- Support MKCALENDAR method
- Support REPORT method for queries
- Support PROPFIND for calendar properties
- Implement sync token support
- Handle ICS content type
- Achieve 75%+ test coverage

## Tasks

### 1. WebDAV Handler Setup
- [ ] Review existing `app/fs.go` - WebDAV filesystem implementation
  - [ ] Understand current Dir implementation
  - [ ] Identify where to add CalDAV support
- [ ] Create `app/webdav/handler.go` - Custom WebDAV handler
  - [ ] Wrap standard webdav.Handler
  - [ ] Add CalDAV method handling
  - [ ] Add content-type detection
  - [ ] Add authentication integration
- [ ] Create `app/webdav/auth.go` - WebDAV auth integration
  - [ ] Extract auth from WebDAV request
  - [ ] Integrate with auth manager (Phase 3)
  - [ ] Set user context for filesystem operations

### 2. Calendar Filesystem
- [ ] Create `app/webdav/calendar_fs/` - Calendar filesystem operations
  - [ ] Create `mkcalendar.go` - MKCALENDAR method
    - [ ] Handle MKCOL request for calendar collection
    - [ ] Create calendar in LevelDB
    - [ ] Set C:calendar-collection property
    - [ ] Validate calendar properties
    - [ ] Return 201 Created on success
    - [ ] Return 405 Method Not Allowed if collection exists
    - [ ] Return 403 Forbidden if no permission
  - [ ] Create `propfind.go` - PROPFIND method
    - [ ] Handle PROPFIND for calendar properties
    - [ ] Return D:resourcetype (C:calendar-collection)
    - [ ] Return C:calendar-home-set
    - [ ] Return C:calendar-description
    - [ ] Return C:calendar-timezone
    - [ ] Return C:supported-calendar-component-set
    - [ ] Return D:getcontenttype (text/calendar)
    - [ ] Return D:getlastmodified
    - [ ] Handle depth headers (0, 1, infinity)
  - [ ] Create `report.go` - REPORT method
    - [ ] Handle REPORT for calendar queries
    - [ ] Parse CalDAV query XML
    - [ ] Support calendar-query filter
    - [ ] Support date-range filter
    - [ ] Support comp-filter (VEVENT, VTODO, etc.)
    - [ ] Support prop-filter (property filters)
    - [ ] Return multistatus with matching items
    - [ ] Return expanded recurrence instances
  - [ ] Create `sync.go` - Sync support
    - [ ] Generate sync token on changes
    - [ ] Handle sync-collection-report
    - [ ] Return changed/deleted items since token
    - [ ] Support sync-token header
- [ ] Create `app/webdav/calendar_fs/types.go` - CalDAV types
  - [ ] CalendarCollection struct
  - [ ] CalendarProperties struct
  - [ ] CalDAVQuery struct (for REPORT)
  - [ ] SyncToken struct

### 3. Filesystem Mapping
- [ ] Create `app/webdav/filesystem.go` - Path to LevelDB mapping
  - [ ] Map `/webdav/{user}/calendars/{cal_uid}/` → Calendar
  - [ ] Map `/webdav/{user}/calendars/{cal_uid}/events.ics` → Events in calendar
  - [ ] Map `/webdav/{user}/calendars/{cal_uid}/todos.ics` → Todos in calendar
  - [ ] Map `/webdav/{user}/calendars/{cal_uid}/journals.ics` → Journals in calendar
  - [ ] Map `/webdav/public/{hash}/` → Public calendar
  - [ ] Resolve symlinks for shared calendars
- [ ] Create `app/webdav/resolver.go` - Path resolution
  - [ ] Resolve user from auth context
  - [ ] Resolve calendar UID from path
  - [ ] Resolve public calendar from hash
  - [ ] Check access permissions
  - [ ] Return resolved path or error

### 4. Content Type Handling
- [ ] Create `app/webdav/content_type.go` - Content type detection
  - [ ] Detect `text/calendar` content type
  - [ ] Set Content-Type header for ICS files
  - [ ] Validate ICS content on upload
  - [ ] Return 415 Unsupported Media Type for non-ICS
- [ ] Create `app/webdav/import_handler.go` - ICS upload handling
  - [ ] Handle PUT request for ICS file
  - [ ] Parse uploaded ICS
  - [ ] Store in LevelDB (create/update event)
  - [ ] Handle recurrence (store rule + expand instances)
  - [ ] Return 201 Created or 204 No Content
- [ ] Create `app/webdav/export_handler.go` - ICS download handling
  - [ ] Handle GET request for ICS file
  - [ ] Generate ICS from LevelDB data
  - [ ] Stream response
  - [ ] Set Content-Disposition header
  - [ ] Handle date range queries via query params

### 5. WebDAV Methods Implementation
- [ ] **PROPFIND** - List properties
  - [ ] Implement for calendar collections
  - [ ] Implement for calendar items (events, etc.)
  - [ ] Support property names (D:, C:)
  - [ ] Return proper multistatus response
- [ ] **PROPPATCH** - Modify properties
  - [ ] Implement for calendar properties
  - [ ] Allow updating calendar description, color, etc.
- [ ] **MKCOL** - Create collection
  - [ ] Implement MKCALENDAR (MKCOL with calendar properties)
  - [ ] Validate calendar properties
  - [ ] Create calendar in LevelDB
- [ ] **REPORT** - Query calendar
  - [ ] Parse CalDAV query XML
  - [ ] Execute query against LevelDB
  - [ ] Return multistatus with results
  - [ ] Handle free/busy queries
- [ ] **GET** - Download item
  - [ ] Handle ICS file download
  - [ ] Return proper content-type
- [ ] **PUT** - Upload item
  - [ ] Handle ICS file upload
  - [ ] Parse and store
- [ ] **DELETE** - Delete item
  - [ ] Handle event/todo deletion
  - [ ] Handle recurrence deletion
- [ ] **COPY** - Copy item
  - [ ] Handle event copy (create new UID)
- [ ] **MOVE** - Move item
  - [ ] Handle event move between calendars

### 6. CalDAV Properties
- [ ] Implement C:calendar-collection property
- [ ] Implement C:calendar-home-set property
- [ ] Implement C:calendar-description property
- [ ] Implement C:calendar-timezone property
- [ ] Implement C:calendar-color property
- [ ] Implement C:calendar-order property
- [ ] Implement C:calendar-timezone-url property
- [ ] Implement C:supported-calendar-component-set property
- [ ] Implement C:calendar-tag property
- [ ] Implement C:calendar-supported-report-set property

### 7. CalDAV Queries
- [ ] Implement calendar-query (REPORT)
  - [ ] comp-filter for VEVENT, VTODO, VJOURNAL
  - [ ] prop-filter for property filters
  - [ ] param-filter for parameter filters
  - [ ] time-range filter
  - [ ] match types (any-of, all-of)
- [ ] Implement freebusy-query (REPORT)
  - [ ] Request free/busy for date range
  - [ ] Aggregate across calendars
  - [ ] Return VFREEBUSY
- [ ] Implement sync-collection-report
  - [ ] Return changes since sync-token
  - [ ] Include added, modified, deleted items

### 8. Sync Tokens
- [ ] Implement sync token generation
  - [ ] Generate on calendar creation
  - [ ] Generate on event creation/update
  - [ ] Include timestamp + version
- [ ] Implement sync token validation
  - [ ] Verify token format
  - [ ] Check token validity
- [ ] Implement sync state tracking
  - [ ] Store sync state per user/calendar
  - [ ] Track last sync token

### 9. WebDAV Tests
- [ ] Create `app/webdav/handler_test.go`
  - [ ] Test WebDAV method routing
  - [ ] Test auth integration
- [ ] Create `app/webdav/calendar_fs/mkcalendar_test.go`
  - [ ] Test MKCALENDAR success
  - [ ] Test MKCALENDAR failure (exists, no permission)
- [ ] Create `app/webdav/calendar_fs/propfind_test.go`
  - [ ] Test calendar property retrieval
  - [ ] Test item property retrieval
- [ ] Create `app/webdav/calendar_fs/report_test.go`
  - [ ] Test calendar-query
  - [ ] Test date-range query
  - [ ] Test recurrence expansion
- [ ] Create `app/webdav/calendar_fs/sync_test.go`
  - [ ] Test sync token generation
  - [ ] Test sync-collection-report
- [ ] Create `app/webdav/import_handler_test.go`
  - [ ] Test ICS upload
  - [ ] Test ICS parsing
- [ ] Create `app/webdav/export_handler_test.go`
  - [ ] Test ICS download
  - [ ] Test content-type header
- [ ] Create `app/webdav/integration_test.go`
  - [ ] Test full WebDAV workflow
  - [ ] Test with actual WebDAV client (if available)
  - [ ] Test RFC 4791 compliance

## WebDAV Path Structure

```
/webdav/
  /{user}/
    /calendars/
      /{cal_uid}/              # Calendar collection
        # PROPFIND returns calendar properties
        # MKCOL creates calendar
        /events.ics            # All events in calendar
        /todos.ics             # All todos in calendar
        /journals.ics          # All journals in calendar
        /{event_uid}.ics       # Single event file
      /{public_hash}/          # Public calendar (no auth)
  /public/
    /{hash}/                   # Alternative public calendar path
```

## Sample WebDAV Requests

### MKCALENDAR - Create Calendar
```xml
MKCOL /webdav/user1/calendars/work/
Calendar-Major: 1.0

Response: 201 Created
C:calendar-home-set: /webdav/user1/calendars/
D:resourcetype: <C:calendar-collection/>
```

### PROPFIND - Get Calendar Properties
```xml
PROPFIND /webdav/user1/calendars/work/
Depth: 0
Prop:
  <D:resourcetype/>
  <C:calendar-description/>
  <C:calendar-color/>

Response: 207 Multi-Status
  <response>
    <href>/webdav/user1/calendars/work/</href>
    <propstat>
      <prop>
        <D:resourcetype><C:calendar-collection/></D:resourcetype>
        <C:calendar-description>Work calendar</C:calendar-description>
        <C:calendar-color>#FF0000</C:calendar-color>
      </prop>
      <status>HTTP/1.1 200 OK</status>
    </propstat>
  </response>
```

### REPORT - Query Events
```xml
REPORT /webdav/user1/calendars/work/
Depth: 0
CalDAV:calendar-query
  <C:comp-filter name="VEVENT">
    <C:time-range start="20260301T000000Z" end="20260331T235959Z"/>
  </C:comp-filter>
  <D:prop>
    <D:displayname/>
    <C:calendar-data/>
  </D:prop>

Response: 207 Multi-Status with matching events
```

### Sync-Collection-Report
```xml
REPORT /webdav/user1/calendars/work/
Depth: 0
Sync-token: abc123
C:sync-collection

Response: 207 Multi-Status with changed items
```

## Success Criteria
- [ ] MKCALENDAR method works (create calendar)
- [ ] PROPFIND returns correct calendar properties
- [ ] REPORT query works (date range, comp-filter)
- [ ] Sync tokens work (generate, validate, sync-report)
- [ ] ICS upload (PUT) works
- [ ] ICS download (GET) works with correct content-type
- [ ] Public calendar access works (via hash)
- [ ] Shared calendar access works (via symlink resolution)
- [ ] Auth integration works (WebDAV + JWT/Basic)
- [ ] 75%+ test coverage
- [ ] All tests pass with `go test -race ./app/webdav/...`

## Notes
- RFC 4791 is complex - test with real CalDAV clients if possible
- CalDAV clients may expect specific behavior - be flexible
- Sync tokens should be efficient (don't store full history)
- REPORT queries should be optimized (use LevelDB indexes)
- Consider implementing partial support if full RFC 4791 is too complex

## Next Phase
Move to **Phase 7: Admin & IANA Timezone Integration** when all items above are complete.
