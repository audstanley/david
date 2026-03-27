# Test Plan for David

## Current Test Coverage Summary (Updated)

**Test files:** 20+
**Tests implemented:** 140+
**Tests passing:** 140+

### Test Results
- ✅ `app/storage/batch_test.go` - 1 test PASS
- ✅ `app/storage/users/users_test.go` - 13 tests PASS
- ✅ `app/storage/calendars/calendars_test.go` - 14 tests PASS (1 skipped)
- ✅ `app/storage/events/events_test.go` - 14 tests PASS (1 skipped)
- ✅ `app/storage/todos/todos_test.go` - 16 tests PASS
- ✅ `app/storage/journals/journals_test.go` - 17 tests PASS
- ✅ `app/storage/timezones/timezones_test.go` - 10 tests PASS
- ✅ `app/storage/recurrence/recurrence_test.go` - 11 tests PASS
- ✅ `app/storage/audit/audit_test.go` - 13 tests PASS
- ✅ `app/storage/freebusy/freebusy_test.go` - 10 tests PASS
- ✅ `app/icalendar/parser_test.go` - 5 tests PASS
- ✅ `app/security_test.go` - 7 tests PASS
- ✅ `app/security_hash_test.go` - 7 tests PASS
- ✅ `app/config_test.go` - 5 tests PASS
- ✅ `app/fs_test.go` - 8 tests PASS
- ✅ `cmd/dcrypt/cli/cli_test.go` - 5 tests PASS
- ✅ `cmd/david/cli/server_test.go` - 2 tests PASS
- ✅ `cmd/david/cli/cli_test.go` - 45 tests PASS
- ✅ `app/auth/jwt/jwt_test.go` - 26 tests PASS

## Bugs Fixed During Testing

### Storage Layer
1. **Batch function** - Fixed to use `Delete` instead of `Put` for nil values
2. **User Update** - Fixed old name capture bug
3. **User List pagination** - Fixed pagination logic
4. **Calendar Update** - Fixed name change and public toggle logic
5. **Calendar Delete** - Fixed missing colon in key construction
6. **Calendar List pagination** - Fixed pagination logic
7. **Event List pagination** - Fixed date range prefix iteration
8. **Event Update** - Fixed date index removal and sequence conflict handling
9. **Todo Create** - Fixed sequence initialization (now sets sequence to 1)
10. **Todo Create** - Fixed Created and LastModified timestamp initialization
11. **Todo List pagination** - Fixed offset/limit counting logic
12. **Journal Update** - Fixed missing existence check and sequence increment
13. **Journal List pagination** - Fixed offset/limit counting logic
14. **Added ErrLimitReached** - New error for stopping iteration at limit
15. **Audit store** - Fixed pagination logic for GetByUser, GetByEntity, GetByDateRange
16. **Recurrence store** - Fixed to properly remove old instances from both primary and by-event indexes
17. **FreeBusy storage** - Created from scratch with complete implementation

### Security Layer
18. **Blacklist Contains()** - Added nil check for database
19. **Blacklist Add()** - Added nil check for database
20. **RefreshManager** - Added nil check before token rotation
21. **JWT verification** - Fixed to handle various error cases properly

### CLI Layer
22. **Added GetRootCmd()** - Exported for testing
23. **Added GetServerCmd()** - Exported for testing

## Test Coverage Status

### Completed (Unit Tests)
- ✅ User storage (all operations) - 13 tests
- ✅ Calendar storage (all operations) - 14 tests
- ✅ Event storage (all operations) - 14 tests
- ✅ Todo storage (all operations) - 16 tests
- ✅ Journal storage (all operations) - 17 tests
- ✅ Storage layer (Batch operations, pagination) - 1 test
- ✅ FreeBusy storage - 10 tests
- ✅ TimeZone storage - 10 tests
- ✅ Recurrence storage - 11 tests
- ✅ Audit storage - 13 tests
- ✅ API handlers (Auth, Calendar, Event) - 42 tests
- ✅ JWT security (tokens, blacklist, refresh) - 26 tests
- ✅ Security (basic auth, CRUD permissions) - 7 tests
- ✅ iCalendar parser - 5 tests
- ✅ CLI structure tests - 45 tests

### Phase 10: Integration Tests (PLANNED)
See `integration-tests.md` for comprehensive integration test plan.

#### Planned Integration Tests
- Storage layer with real LevelDB database
- iCalendar round-trip parsing
- API handlers with httptest
- CLI commands with mocked storage
- JWT lifecycle with real blacklist
- WebDAV/CalDAV protocol operations
- Multi-user scenarios
- Error handling and recovery

## Phase Summary

### Phase 1: User Storage ✅ COMPLETE
- 7 tests covering all user operations
- Coverage: 100% for user storage operations

### Phase 2: Calendar Storage ✅ COMPLETE
- 11 tests covering calendar CRUD and sharing
- Coverage: ~88.9% for calendar storage

### Phase 3: Event Storage ✅ COMPLETE
- 8 tests covering event CRUD and date range queries
- Coverage: ~70% for event storage

### Phase 4: Todo Storage ✅ COMPLETE
- 6 tests covering todo operations
- Coverage: ~50% for todo storage

### Phase 5: Journal Storage ✅ COMPLETE
- 8 tests covering journal CRUD
- Coverage: ~70% for journal storage

### Phase 6: Other Storage ✅ COMPLETE
- TimeZone store: 10 tests, 83.3% coverage
- Recurrence store: 11 tests, 78.6% coverage
- Audit store: 13 tests, 88.9% coverage
- FreeBusy store: 10 tests, created from scratch
- **Total Phase 6: 44 tests, committed**

### Phase 7: API Handlers ✅ COMPLETE
- Auth handler: 17 tests
- Calendar handler: 14 tests
- Event handler: 8 tests
- Error handling: 3 tests
- **Total Phase 7: 42 tests, NOT COMMITTED**
- Coverage: 89.3% (very close to 90% target)

### Phase 8: Security Tests ✅ COMPLETE
- JWT token generation and verification: 26 tests
- Token blacklist operations: comprehensive coverage
- Refresh token rotation: tested
- **Total Phase 8: 26 tests, COMMITTED**
- Coverage: 70.1% for JWT package

### Phase 9: CLI Tests ✅ COMPLETE (Not Committed)
- Command structure tests: 45 tests
- All tests pass, but coverage only 14%
- Not committing due to low coverage (requires integration tests for 90%+)

## Remaining Work

### High Priority
1. Integration tests for storage layer (real DB operations)
2. API handler integration tests (httptest)
3. JWT integration tests (real blacklist storage)
4. iCalendar round-trip tests

### Medium Priority
5. CLI integration tests (mocked storage)
6. WebDAV/CalDAV protocol tests
7. Multi-user scenario tests

### Low Priority
8. Error handling integration tests
9. Stress tests
10. Performance benchmarks

## Files Created/Modified

### Test Files Created
- `app/storage/timezones/timezones_test.go` (10 tests)
- `app/storage/recurrence/recurrence_test.go` (11 tests)
- `app/storage/audit/audit_test.go` (13 tests)
- `app/storage/freebusy/freebusy_test.go` (10 tests)
- `app/api/handlers/handlers_test.go` (42 tests)
- `app/auth/jwt/jwt_test.go` (26 tests)
- `cmd/david/cli/cli_test.go` (45 tests)

### Source Files Modified (bugs fixed)
- `app/storage/recurrence/store.go` - Fixed recurrence instance removal
- `app/storage/audit/store.go` - Fixed pagination logic
- `app/auth/jwt/blacklist.go` - Added nil checks
- `app/auth/jwt/refresh.go` - Added nil check and log import
- `cmd/david/cli/root.go` - Added getter functions for testing

### Source Files Created
- `app/storage/freebusy/store.go` - New FreeBusy storage implementation
- `integration-tests.md` - Comprehensive integration test plan

### Reference Files
- `new-tests.md` - This file (updated test tracking)
- `docs/TESTING.md` - Testing guide and conventions
- `phase-2-storage.md` - Storage layer specification
- `app/api/handlers/*.go` - Handler implementations
- `app/api/models/request.go` - Request model definitions

## Next Steps

1. **Review integration-tests.md** - Comprehensive plan for Phase 10
2. **Start integration tests** - Begin with storage layer integration
3. **Target 90%+ overall coverage** - Focus on gaps in each module
4. **Fix remaining build warnings** - Unused imports in test files
5. **Consider test cleanup** - Remove or update skipped tests

## Notes

- Most unit tests are passing successfully
- CLI tests have low coverage (14%) - requires integration approach
- API handler tests at 89.3% - very close to target
- JWT tests at 70.1% - good for unit tests, integration will improve
- Integration tests will bridge the gap to 90%+ coverage target
