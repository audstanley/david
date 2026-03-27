# Test Plan for David

## Current Test Coverage Summary (Updated)

**Test files:** 16
**Tests implemented:** 74+
**Tests passing:** 73+
**Tests skipped:** 2 (test framework issues)

### Test Results
- ✅ `app/storage/batch_test.go` - 1 test PASS
- ✅ `app/storage/users/users_test.go` - 13 tests PASS
- ✅ `app/storage/calendars/calendars_test.go` - 14 tests PASS (1 skipped)
- ✅ `app/storage/events/events_test.go` - 14 tests PASS (1 skipped)
- ✅ `app/storage/todos/todos_test.go` - 16 tests PASS
- ✅ `app/storage/journals/journals_test.go` - 17 tests PASS
- ✅ `app/icalendar/parser_test.go` - 5 tests PASS
- ✅ `app/security_test.go` - 7 tests PASS
- ✅ `app/security_hash_test.go` - 7 tests PASS
- ✅ `app/config_test.go` - 5 tests PASS
- ✅ `app/fs_test.go` - 8 tests PASS
- ✅ `cmd/dcrypt/cli/cli_test.go` - 5 tests PASS
- ✅ `cmd/david/cli/server_test.go` - 2 tests PASS
- ⚠️ `cmd/david/cli/cli_test.go` - 4 tests (builds with errors)

## Bugs Fixed During Testing

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

## Test Coverage Status

### Completed
- ✅ User storage (all operations)
- ✅ Calendar storage (all operations)
- ✅ Event storage (all operations)
- ✅ Todo storage (all operations)
- ✅ Storage layer (Batch delete fix, pagination)

### Pending
- ⏳ Journal storage
- ⏳ API handlers
- ⏳ Security tests (JWT, RBAC)
- ⏳ CLI integration tests

## Next Steps

1. Implement Event storage tests
2. Implement Todo storage tests
3. Implement Journal storage tests
4. Implement API handler tests
5. Fix CLI tests build errors

**10 test files** with ~45-50 individual tests:
- `app/storage/batch_test.go` - 1 test
- `app/storage/users/users_test.go` - 1 test
- `app/icalendar/parser_test.go` - 5 tests
- `app/security_test.go` - 7 tests
- `app/security_hash_test.go` - 7 tests
- `app/config_test.go` - 5 tests
- `app/fs_test.go` - 8 tests
- `cmd/dcrypt/cli/cli_test.go` - 5 tests
- `cmd/david/cli/server_test.go` - 2 tests
- `cmd/david/cli/cli_test.go` - 4 tests (builds with errors)

## Missing Test Coverage

### Critical Areas Without Tests
1. User storage (Create, Get, Update, List - only Delete tested)
2. Calendar storage (all operations)
3. Event storage (all operations)
4. Todo storage (all operations)
5. Journal storage (all operations)
6. FreeBusy storage (all operations)
7. TimeZone storage (all operations)
8. Recurrence storage (all operations)
9. Audit storage (all operations)
10. API handlers (auth, calendar, event)
11. JWT token operations
12. RBAC authorization
13. Basic auth handler

### High Priority Tests

#### User Storage Tests
- Create user with valid data
- Create user with duplicate username
- Get user by ID
- Get user by username
- Update user fields (username, email, display name, role)
- Update user password
- List users with pagination
- Delete user (already tested)

#### Calendar Storage Tests
- Create calendar
- Get calendar by UID
- Get calendar by name
- Update calendar
- Delete calendar
- List calendars with pagination
- Check calendar is public
- Share calendar with user
- Check access permissions

#### Event Storage Tests
- Create event
- Get event by UID
- Update event (with sequence number)
- Delete event
- List events in date range
- Get event instance in date range
- Create recurring event
- Get recurrence instances
- Delete recurrence

#### Todo Storage Tests
- Create todo
- Get todo by UID
- Update todo
- Delete todo
- List todos

#### Journal Storage Tests
- Create journal
- Get journal by UID
- Update journal
- Delete journal
- List journals

#### API Handler Tests
- Auth handler: login, logout, token refresh
- Calendar handler: CRUD operations
- Event handler: CRUD operations
- Error handling and validation

#### Security Tests
- JWT token generation and validation
- JWT blacklist operations
- RBAC role checks
- Basic auth verification

### Medium Priority Tests
- Integration tests for CLI commands
- Error handling in storage layer
- Edge cases and boundary conditions
- Concurrent access tests

### Low Priority Tests
- Backup and restore operations
- Database compaction
- Statistics retrieval
- Timezone resolution

---

## TODO Lists

### Phase 1: User Storage (High Priority) - COMPLETE
- [x] Test user creation with valid data
- [x] Test user creation with duplicate username
- [x] Test GetByID for non-existent user
- [x] Test GetByUsername
- [x] Test Update user fields
- [x] Test Update password
- [x] Test List with pagination
- [x] Test List with offset/limit

### Phase 2: Calendar Storage (High Priority) - COMPLETE
- [x] Test calendar creation
- [x] Test GetByUID
- [x] Test GetByName
- [x] Test Update calendar
- [x] Test Delete calendar
- [x] Test List with pagination
- [x] Test IsPublic

### Phase 3: Event Storage (High Priority) - COMPLETE
- [x] Test event creation
- [x] Test GetByUID
- [x] Test Update event
- [x] Test Delete event
- [x] Test List events in date range
- [x] Test List pagination
- [x] Test GetInDateRange
- [x] Test sequence conflict check

### Phase 4: Todo Storage (Medium Priority) - COMPLETE
- [x] Test todo creation
- [x] Test GetByUID
- [x] Test Update todo
- [x] Test Delete todo
- [x] Test List todos with pagination

### Phase 5: Journal Storage (Medium Priority) - COMPLETE
- [x] Test journal creation
- [x] Test GetByUID
- [x] Test Update journal
- [x] Test Delete journal
- [x] Test List journals with pagination

### Remaining Phases
- [ ] Phase 6: Other Storage (FreeBusy, TimeZone, Recurrence, Audit)
- [ ] Phase 7: API Handlers
- [ ] Phase 8: Security Tests
- [ ] Phase 9: CLI Tests

## Phase 2: Calendar Storage (High Priority) - ALMOST COMPLETE
- [x] Test calendar creation
- [x] Test GetByUID
- [x] Test GetByName
- [x] Test Update calendar
- [x] Test Delete calendar
- [x] Test List with pagination
- [x] Test IsPublic
- [x] Test calendar share creation
- [x] Test calendar share deletion
- [x] Test CheckAccess permissions
- [ ] Test calendar duplicate UID (test framework issue, works in isolation)

### Phase 2: Calendar Storage (High Priority)
- [ ] Test calendar creation
- [ ] Test GetByUID
- [ ] Test GetByName
- [ ] Test Update calendar
- [ ] Test Delete calendar
- [ ] Test List with pagination
- [ ] Test IsPublic
- [ ] Test calendar share creation
- [ ] Test calendar share deletion
- [ ] Test CheckAccess permissions

### Phase 3: Event Storage (High Priority)
- [ ] Test event creation
- [ ] Test GetByUID
- [ ] Test Update event
- [ ] Test Delete event
- [ ] Test List events in date range
- [ ] Test GetInDateRange
- [ ] Test Create recurring event
- [ ] Test Get recurrence instances
- [ ] Test Delete recurrence

### Phase 4: Todo Storage (Medium Priority) - COMPLETE
- [x] Test todo creation
- [x] Test GetByUID
- [x] Test Update todo
- [x] Test Delete todo
- [x] Test List todos with pagination

### Phase 5: Journal Storage (Medium Priority) - COMPLETE
- [x] Test journal creation
- [x] Test GetByUID
- [x] Test Update journal
- [x] Test Delete journal
- [x] Test List journals with pagination

### Phase 6: Other Storage (Medium Priority)
- [ ] FreeBusy generation
- [ ] TimeZone store
- [ ] Recurrence store
- [ ] Audit log

### Phase 7: API Handlers (High Priority)
- [ ] Auth handler tests
- [ ] Calendar handler tests
- [ ] Event handler tests
- [ ] Error handling

### Phase 8: Security (Medium Priority)
- [ ] JWT token tests
- [ ] RBAC role tests
- [ ] Basic auth tests

### Phase 9: CLI Tests (Low Priority)
- [ ] Fix broken CLI tests
- [ ] User command tests
- [ ] Calendar command tests
- [ ] Event command tests

---

## Implementation Priority

**Start with Phase 1 (User Storage)** - Complete the test suite for user operations since this is the area with the bug that was just fixed.
