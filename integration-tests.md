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

## Notes

- Integration tests may take longer than unit tests
- Use `-short` flag for CI to skip expensive tests
- Consider parallel test execution
- Use testcontainers for external dependencies if needed
- Document any test-only code or mocks
