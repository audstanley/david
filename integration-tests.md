# Integration Test Plan for David

## Overview
This document outlines the comprehensive integration test strategy for the David CalDAV server project. Integration tests verify that different components work together correctly.

## Test Coverage Goals

### Phase 10: Integration Tests (90%+ Coverage Target)

#### 10.1 Storage Layer Integration (Target: 85%)
- Test actual LevelDB storage operations with real database
- Test cross-entity relationships (events linked to calendars)
- Test batch operations
- Test pagination with real data
- Test concurrent access

#### 10.2 iCalendar Parser Integration (Target: 90%)
- Parse real iCalendar files (RFC 2445 compliant)
- Round-trip test: parse → generate → parse
- Test with various calendar components:
  - Events with recurrence
  - Todos with due dates
  - Journals with status
  - FreeBusy blocks
  - VTIMEZONE components

#### 10.3 API Handler Integration (Target: 95%)
- Test HTTP handlers with httptest
- Test authentication middleware
- Test authorization (RBAC)
- Test error handling
- Test request/response validation
- Test all endpoints:
  - Auth endpoints (login, logout, refresh)
  - Calendar CRUD
  - Event CRUD
  - Todo CRUD
  - Journal CRUD
  - User management

#### 10.4 CLI Integration (Target: 80%)
- Test CLI commands with mocked storage
- Test import/export functionality
- Test user management commands
- Test calendar management commands
- Test event management commands
- Test admin commands

#### 10.5 Security Integration (Target: 90%)
- JWT token lifecycle tests
- Token blacklist operations with real storage
- Refresh token rotation
- Concurrent token operations
- Basic auth with real users
- RBAC authorization checks

#### 10.6 WebDAV/CalDAV Protocol Integration (Target: 85%)
- PROPFIND operations
- PROPPATCH operations
- MKCOL operations
- COPY/MOVE operations
- LOCK/UNLOCK operations
- REPORT operations
- Calendar collection operations

#### 10.7 Multi-User Scenarios (Target: 85%)
- Multiple users accessing same calendar
- Calendar sharing and access control
- Concurrent user operations
- Conflict resolution

#### 10.8 Error Handling Integration (Target: 90%)
- Database errors
- Network errors
- Permission errors
- Invalid input handling
- Recovery from panics

## Test Architecture

### Test Helpers
```
app/test/
  fixtures/          # Test data files
  mocks/             # Mock implementations
    storage_mock.go  # Mock storage layer
    db_mock.go       # Mock database
  helpers.go         # Common test utilities
```

### Test Patterns

#### Unit + Integration Hybrid
```go
func TestStorageWithRealDB(t *testing.T) {
    db := storage.NewStorage(t.TempDir())
    defer db.Close()
    
    // Test actual database operations
    // ...
}
```

#### Mock-Based Integration
```go
func TestAPIWithMockStorage(t *testing.T) {
    mockStorage := storage.NewMockStorage()
    handler := NewCalendarHandler(mockStorage)
    
    // Test API layer with controlled data
    // ...
}
```

#### End-to-End Tests
```go
func TestFullCalendarWorkflow(t *testing.T) {
    // Start server
    // Create user
    // Create calendar
    // Add event
    // Query events
    // Verify data consistency
    // Cleanup
}
```

## Test Data Fixtures

### Sample iCalendar Files
- basic_event.ics
- recurring_event.ics
- recurring_event_instances.ics
- todo.ics
- journal.ics
- busy_block.ics
- vtimezone.ics

### Sample Users
- admin user
- regular user
- reader user

### Sample Calendars
- personal calendar
- work calendar
- shared calendar

## Implementation Priority

1. **High Priority** (Core functionality)
   - Storage layer integration
   - iCalendar parser integration
   - API handler integration
   - Security integration

2. **Medium Priority** (Features)
   - CLI integration
   - WebDAV/CalDAV protocol integration
   - Multi-user scenarios

3. **Low Priority** (Edge cases)
   - Error handling integration
   - Stress tests
   - Performance tests

## Running Tests

```bash
# All tests
go test ./... -v -race

# Coverage
go test ./... -coverprofile=cover.out
go tool cover -html=cover.out

# Specific package
go test ./app/auth/jwt/... -v

# Integration tests only
go test ./app/integration/... -v

# With race detector
go test -race ./...

# Benchmarks
go test -bench=. ./...
```

## Test Environment Setup

### Database
- Use temporary directories for test databases
- Clean up after tests
- Support both fresh and existing databases

### Configuration
- Use test configuration files
- Override default settings
- Mock external dependencies

### Network
- Use localhost for test servers
- Avoid port conflicts
- Clean up server processes

## Expected Coverage by Package

| Package | Current | Target | Strategy |
|---------|---------|--------|----------|
| app/storage | ~60% | 85% | Real DB tests |
| app/icalendar | ~70% | 90% | Parse/generate round-trip |
| app/api/handlers | ~89% | 95% | httptest integration |
| app/auth/jwt | 70% | 90% | Token lifecycle tests |
| cmd/david/cli | 14% | 80% | Mocked storage tests |
| app/security | ~60% | 85% | Auth/authorization tests |
| app/auth/rbac | 0% | 90% | RBAC scenario tests |
| app/auth/basic | 0% | 85% | Basic auth tests |

## Test Maintenance

### Guidelines
- Keep tests isolated and independent
- Use table-driven tests where possible
- Document test scenarios clearly
- Update tests when code changes
- Avoid flaky tests (deterministic behavior)

### CI/CD Integration
- Run all tests on push
- Run integration tests on main branch
- Coverage reporting in PRs
- Block merge if coverage drops below threshold

## TODOs - Integration Test Implementation

### Phase 10.1: Storage Layer Integration (Target: 85%)

#### 10.1.1 LevelDB Integration Tests
- [x] Create test storage helper
- [x] Test create operations with real DB
- [x] Test get operations with real DB
- [x] Test update operations with real DB
- [x] Test delete operations with real DB
- [x] Test list operations with pagination
- [ ] Test batch operations
- [ ] Test concurrent access
- [ ] Test cross-entity relationships

#### 10.1.2 User Storage Integration
- [x] Test user creation with real database
- [x] Test user get by ID
- [x] Test user get by username
- [x] Test user update
- [x] Test user list with pagination
- [ ] Test user deletion cascade

#### 10.1.3 Calendar Storage Integration
- [ ] Test calendar CRUD with real DB
- [ ] Test calendar share operations
- [ ] Test calendar access checks
- [ ] Test calendar list with pagination

#### 10.1.4 Event Storage Integration
- [ ] Test event CRUD with real DB
- [ ] Test event date range queries
- [ ] Test event recurrence storage
- [ ] Test event sequence handling

#### 10.1.5 Todo Storage Integration
- [ ] Test todo CRUD with real DB
- [ ] Test todo list with pagination

#### 10.1.6 Journal Storage Integration
- [ ] Test journal CRUD with real DB
- [ ] Test journal list with pagination

#### 10.1.7 FreeBusy Storage Integration
- [ ] Test freebusy generation
- [ ] Test freebusy storage and retrieval

#### 10.1.8 TimeZone Storage Integration
- [ ] Test timezone store
- [ ] Test timezone retrieval

#### 10.1.9 Recurrence Storage Integration
- [ ] Test recurrence instance storage
- [ ] Test recurrence instance retrieval

#### 10.1.10 Audit Storage Integration
- [ ] Test audit log creation
- [ ] Test audit log queries

### Phase 10.2: iCalendar Parser Integration (Target: 90%)

#### 10.2.1 Basic Parsing Tests
- [x] Parse basic event ICS
- [x] Parse recurring event ICS
- [x] Parse TODO ICS
- [x] Parse JOURNAL ICS
- [x] Parse FREEBUSY ICS

#### 10.2.2 Round-Trip Tests
- [x] Parse → Generate → Parse cycle for events
- [ ] Parse → Generate → Parse cycle for todos
- [ ] Parse → Generate → Parse cycle for journals
- [x] Verify UID consistency through round-trip
- [ ] Verify sequence number handling

#### 10.2.3 Recurrence Tests
- [x] Parse recurrence rules
- [ ] Generate recurrence instances
- [x] Test EXDATE handling
- [ ] Test RDATE handling

#### 10.2.4 Timezone Tests
- [x] Parse VTIMEZONE components
- [ ] Test timezone conversion
- [ ] Test daylight saving time handling

#### 10.2.5 Edge Cases
- [x] Parse invalid ICS (error handling)
- [x] Parse empty calendar
- [x] Parse calendar with multiple components

### Phase 10.3: API Handler Integration (Target: 95%)

#### 10.3.1 Test Setup
- [x] Create test server with httptest
- [x] Create mock storage for API tests
- [x] Create test user fixtures
- [x] Create test calendar fixtures

#### 10.3.2 Auth Handler Integration
- [x] Test login endpoint
- [x] Test token refresh endpoint
- [x] Test logout/blacklist endpoint
- [x] Test invalid credentials handling
- [x] Test token expiration handling

#### 10.3.3 Calendar Handler Integration
- [x] Test calendar creation
- [x] Test calendar retrieval
- [x] Test calendar update
- [x] Test calendar deletion
- [x] Test calendar listing
- [x] Test calendar sharing
- [ ] Test calendar access control

#### 10.3.4 Event Handler Integration
- [x] Test event creation
- [x] Test event retrieval
- [x] Test event update
- [x] Test event deletion
- [x] Test event listing
- [ ] Test date range queries

#### 10.3.5 Todo Handler Integration (Not implemented)
- [ ] Test todo creation
- [ ] Test todo retrieval
- [ ] Test todo update
- [ ] Test todo deletion
- [ ] Test todo listing

#### 10.3.6 Journal Handler Integration (Not implemented)
- [ ] Test journal creation
- [ ] Test journal retrieval
- [ ] Test journal update
- [ ] Test journal deletion
- [ ] Test journal listing

#### 10.3.7 User Handler Integration (Not implemented)
- [ ] Test user creation
- [ ] Test user retrieval
- [ ] Test user update
- [ ] Test user deletion
- [ ] Test user listing

#### 10.3.8 Error Handling (Not implemented)
- [ ] Test validation errors
- [ ] Test not found errors
- [ ] Test permission errors
- [ ] Test server errors

### Phase 10.4: CLI Integration (Target: 80%)

#### 10.4.1 Test Setup
- [ ] Create CLI test helper
- [ ] Create mock storage for CLI
- [ ] Create test database setup

#### 10.4.2 Import/Export Tests
- [ ] Test import single ICS file
- [ ] Test import directory
- [ ] Test export single event
- [ ] Test export calendar
- [ ] Test export date range

#### 10.4.3 User Commands
- [ ] Test user list command
- [ ] Test user create command
- [ ] Test user delete command

#### 10.4.4 Calendar Commands
- [ ] Test calendar list command
- [ ] Test calendar show command
- [ ] Test calendar create command
- [ ] Test calendar delete command
- [ ] Test calendar share command

#### 10.4.5 Event Commands
- [ ] Test event list command
- [ ] Test event show command
- [ ] Test event delete command

#### 10.4.6 Admin Commands
- [ ] Test admin stats command
- [ ] Test admin audit command

### Phase 10.5: Security Integration (Target: 90%)

#### 10.5.1 JWT with Real Storage
- [x] Test token generation with blacklist DB
- [x] Test token verification with blacklist
- [x] Test token revocation
- [x] Test blacklist persistence

#### 10.5.2 Refresh Token Integration
- [x] Test refresh token rotation with DB
- [x] Test old token revocation
- [x] Test refresh token expiration

#### 10.5.3 Basic Auth Integration
- [x] Test basic auth with real users
- [x] Test basic auth permissions
- [x] Test basic auth failures

#### 10.5.4 RBAC Integration
- [x] Test role-based permissions
- [x] Test permission checks in handlers
- [x] Test unauthorized access blocking

### Phase 10.6: WebDAV/CalDAV Protocol Integration (Target: 85%)

#### 10.6.1 Basic Operations
- [ ] Test PROPFIND on root
- [ ] Test PROPFIND on calendar
- [ ] Test PROPFIND on event
- [ ] Test MKCOL for calendar creation
- [ ] Test PUT for event creation

#### 10.6.2 Advanced Operations
- [ ] Test PROPPATCH
- [ ] Test COPY operation
- [ ] Test MOVE operation
- [ ] Test LOCK operation
- [ ] Test UNLOCK operation

#### 10.6.3 CalDAV Specific
- [ ] Test REPORT for calendar queries
- [ ] Test calendar query with date range
- [ ] Test calendar query with UID

### Phase 10.7: Multi-User Scenarios (Target: 85%)

#### 10.7.1 Shared Access
- [ ] Test two users accessing shared calendar
- [ ] Test read-only access
- [ ] Test write access
- [ ] Test admin access

#### 10.7.2 Concurrent Operations
- [ ] Test concurrent event creation
- [ ] Test concurrent event updates
- [ ] Test conflict detection
- [ ] Test sequence number handling

#### 10.7.3 Access Control
- [ ] Test user without access cannot see calendar
- [ ] Test user with read access cannot modify
- [ ] Test user with write access can modify
- [ ] Test calendar owner permissions

### Phase 10.8: Error Handling Integration (Target: 90%)

#### 10.8.1 Database Errors
- [ ] Test disk full scenario
- [ ] Test database corruption handling
- [ ] Test connection failures
- [ ] Test recovery from errors

#### 10.8.2 Network Errors
- [ ] Test timeout handling
- [ ] Test connection reset
- [ ] Test partial data handling

#### 10.8.3 Permission Errors
- [ ] Test unauthorized access
- [ ] Test insufficient permissions
- [ ] Test admin bypass

#### 10.8.4 Input Validation
- [ ] Test invalid email handling
- [ ] Test invalid dates
- [ ] Test invalid ICS content
- [ ] Test SQL injection attempts

#### 10.8.5 Panic Recovery
- [ ] Test panic recovery in handlers
- [ ] Test error logging
- [ ] Test cleanup on panic

## Progress Summary

- Phase 10.1: Storage Layer Integration: 10/35 tasks (29%)
- Phase 10.2: iCalendar Parser Integration: 10/14 tasks (71%)
- Phase 10.3: API Handler Integration: 15/30 tasks (50%)
- Phase 10.4: CLI Integration: 0/15 tasks (0%)
- Phase 10.5: Security Integration: 11/11 tasks (100%)
- Phase 10.6: WebDAV/CalDAV Integration: 0/10 tasks (0%)
- Phase 10.7: Multi-User Scenarios: 0/11 tasks (0%)
- Phase 10.8: Error Handling Integration: 0/15 tasks (0%)

**Total: 46/131 tasks complete**
