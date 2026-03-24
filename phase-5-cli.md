# Phase 5: CLI Commands

## Overview
This phase implements CLI commands for import, export, and calendar/user management. Commands use Cobra for structure and Viper for configuration. This provides a convenient way to manage calendars from the command line.

## Goals
- Implement `david import` command for ICS import
- Implement `david export` command for ICS export
- Implement calendar management commands
- Implement user management commands
- Implement admin commands (statistics, audit)
- Achieve 80%+ test coverage

## Tasks

### 1. CLI Infrastructure
- [ ] Review `cmd/david/cli/root.go` - Root command
  - [ ] Ensure Viper is configured
  - [ ] Set up global flags (config, debug, production)
  - [ ] Add version command
- [ ] Review `cmd/david/cli/server.go` - Server command
  - [ ] Ensure it starts both WebDAV and API server
  - [ ] Add flags for port overrides
- [ ] Create `cmd/david/cli/common.go` - Common utilities
  - [ ] Load config
  - [ ] Connect to LevelDB
  - [ ] Create storage clients
  - [ ] Create API client (for remote operations)

### 2. Import Command
- [ ] Create `cmd/david/cli/import.go` - Import command
  - [ ] Cobra command definition
  - [ ] Flags:
    - `--calendar string` - Import to specific calendar
    - `--dry-run` - Validate without importing
    - `--force` - Overwrite existing events with same UID
    - `--verbose` - Show detailed output
  - [ ] Implement import logic
    - [ ] Parse ICS file(s)
    - [ ] Validate before import
    - [ ] Handle duplicate UIDs
    - [ ] Store in LevelDB
    - [ ] Return summary (imported, skipped, errors)
- [ ] Create `cmd/david/cli/import/validate.go` - Validation
  - [ ] Parse and validate ICS without storing
  - [ ] Check for required properties
  - [ ] Check for invalid values
  - [ ] Return validation errors
- [ ] Create `cmd/david/cli/import/dedupe.go` - Duplicate handling
  - [ ] Check if UID exists
  - [ ] Handle conflicts (skip, overwrite, create new)
  - [ ] Log duplicates
- [ ] Create `cmd/david/cli/import_test.go`
  - [ ] Test import of valid ICS
  - [ ] Test import of multiple files
  - [ ] Test dry-run mode
  - [ ] Test duplicate handling
  - [ ] Test validation errors

### 3. Export Command
- [ ] Create `cmd/david/cli/export.go` - Export command
  - [ ] Cobra command definition
  - [ ] Flags:
    - `--calendar string` - Export specific calendar
    - `--event string` - Export single event by UID
    - `--range string` - Export date range (start,end)
    - `--all` - Export all calendars
    - `--output string` - Output file (default: stdout)
    - `--verbose` - Show detailed output
  - [ ] Implement export logic
    - [ ] Query LevelDB for data
    - [ ] Serialize to ICS
    - [ ] Write to file or stdout
    - [ ] Return summary (exported count)
- [ ] Create `cmd/david/cli/export/filter.go` - Export filtering
  - [ ] Filter events by date range
  - [ ] Filter by calendar
  - [ ] Filter by type (events, todos, journals)
  - [ ] Aggregate multiple calendars
- [ ] Create `cmd/david/cli/export/batch.go` - Batch export
  - [ ] Export multiple calendars
  - [ ] Combine into single ICS file
  - [ ] Handle calendar boundaries
- [ ] Create `cmd/david/cli/export_test.go`
  - [ ] Test calendar export
  - [ ] Test single event export
  - [ ] Test date range export
  - [ ] Test all calendars export
  - [ ] Test round-trip (import → export → compare)

### 4. Calendar Commands
- [ ] Create `cmd/david/cli/calendar/` - Calendar subcommands
  - [ ] Create `list.go` - List calendars
    - [ ] `david calendar list`
    - [ ] Show calendar UID, name, owner, is_public
    - [ ] Show event counts
  - [ ] Create `show.go` - Show calendar details
    - [ ] `david calendar show {uid}`
    - [ ] Show full calendar metadata
    - [ ] Show shares
  - [ ] Create `create.go` - Create calendar
    - [ ] `david calendar create --name "My Calendar"`
    - [ ] Accept --description, --color, --public flags
    - [ ] Generate UID
  - [ ] Create `update.go` - Update calendar
    - [ ] `david calendar update {uid} --name "New Name"`
    - [ ] Accept any mutable properties
  - [ ] Create `delete.go` - Delete calendar
    - [ ] `david calendar delete {uid}`
    - [ ] Confirm deletion
    - [ ] Handle associated events
  - [ ] Create `export.go` - Export calendar
    - [ ] `david calendar export {uid} output.ics`
    - [ ] Wrapper around export command
  - [ ] Create `share.go` - Share calendar
    - [ ] `david calendar share {uid} --user user123 --role write`
    - [ ] Grant access
    - [ ] `david calendar unshare {uid} --user user123`
    - [ ] Revoke access
- [ ] Create `cmd/david/cli/calendar_test.go`
  - [ ] Test calendar CRUD
  - [ ] Test share/unshare

### 5. Event Commands
- [ ] Create `cmd/david/cli/event/` - Event subcommands
  - [ ] Create `list.go` - List events
    - [ ] `david event list --calendar {uid}`
    - [ ] Accept --start, --end date flags
    - [ ] Show event summary (UID, summary, start, end)
  - [ ] Create `show.go` - Show event
    - [ ] `david event show {uid}`
    - [ ] Show full event details
    - [ ] Show expanded recurrence instances
  - [ ] Create `delete.go` - Delete event
    - [ ] `david event delete {uid}`
    - [ ] Confirm deletion
    - [ ] Handle recurrence (all instances or single)
- [ ] Create `cmd/david/cli/event_test.go`
  - [ ] Test event listing
  - [ ] Test event show
  - [ ] Test event deletion

### 6. User Commands
- [ ] Create `cmd/david/cli/user/` - User subcommands
  - [ ] Create `list.go` - List users
    - [ ] `david user list`
    - [ ] Admin only
    - [ ] Show user ID, username, role, email
  - [ ] Create `create.go` - Create user
    - [ ] `david user create --username alice --password secret --role user`
    - [ ] Hash password
    - [ ] Generate user ID
  - [ ] Create `delete.go` - Delete user
    - [ ] `david user delete {user}`
    - [ ] Admin only
    - [ ] Handle user's calendars (reassign or delete)
  - [ ] Create `role.go` - Change user role
    - [ ] `david user role {user} admin`
    - [ ] `david user role {user} user`
    - [ ] Admin only
  - [ ] Create `apikey.go` - API key management
    - [ ] `david user apikey create --name "Desktop"`
    - [ ] Generate API key
    - [ ] Display key (once)
    - [ ] `david user apikey list`
    - [ ] `david user apikey revoke {key_id}`
- [ ] Create `cmd/david/cli/user_test.go`
  - [ ] Test user creation
  - [ ] Test role changes
  - [ ] Test API key management

### 7. Admin Commands
- [ ] Create `cmd/david/cli/admin/` - Admin subcommands
  - [ ] Create `stats.go` - Show statistics
    - [ ] `david admin stats`
    - [ ] Total users, calendars, events, storage size
    - [ ] Per-calendar statistics
    - [ ] Per-user statistics
  - [ ] Create `audit.go` - Show audit log
    - [ ] `david admin audit --user {user}`
    - [ ] `david admin audit --entity {type}:{uid}`
    - [ ] `david admin audit --since {date}`
    - [ ] Show recent audit entries
  - [ ] Create `ratelimit.go` - Rate limit management
    - [ ] `david admin ratelimit list`
    - [ ] `david admin ratelimit add --type user --id {id} --limit 200`
    - [ ] `david admin ratelimit remove {id}`
- [ ] Create `cmd/david/cli/admin_test.go`
  - [ ] Test stats command
  - [ ] Test audit query

### 8. Setup Command
- [ ] Create `cmd/david/cli/setup.go` - Initial setup
  - [ ] `david setup`
  - [ ] Create first admin user
  - [ ] Create initial configuration
  - [ ] Initialize databases
  - [ ] Interactive mode (prompt for password)
  - [ ] Non-interactive mode (use flags)

### 9. Test Fixtures
- [ ] Create `test/fixtures/cli/import/` - Test ICS files for import
  - [ ] `event_basic.ics`
  - [ ] `recurring.ics`
  - [ ] `multi_calendar.ics`
- [ ] Create `test/fixtures/cli/export/` - Expected export outputs
  - [ ] Golden files for export tests

### 10. Unit Tests
- [ ] All test files for each command (listed above)
- [ ] Test CLI flag parsing
- [ ] Test command execution with mocks
- [ ] Test error handling
- [ ] Test with real LevelDB (integration tests)

## Command Examples

```bash
# Import
david import calendar.ics                              # Import to user's default calendar
david import --calendar mycal events.ics               # Import to specific calendar
david import --dir /path/to/*.ics                      # Import multiple files
david import --dry-run events.ics                      # Validate without importing
david import --force overwriting.ics                   # Overwrite existing UIDs

# Export
david export --calendar mycal mycal.ics                # Export calendar
david export --event {uid} event.ics                   # Export single event
david export --range 20260301,20260331 out.ics         # Export date range
david export --all all-calendars.ics                   # Export all calendars

# Calendar
david calendar list                                      # List all calendars
david calendar show {uid}                                # Show calendar details
david calendar create --name "Work Calendar" --public  # Create calendar
david calendar update {uid} --name "Updated Name"      # Update calendar
david calendar delete {uid}                              # Delete calendar
david calendar share {uid} --user user123 --role write  # Share calendar
david calendar unshare {uid} --user user123            # Unshare calendar

# Event
david event list --calendar {uid} --start 20260301     # List events
david event show {uid}                                   # Show event
david event delete {uid}                                 # Delete event

# User
david user list                                          # List users (admin)
david user create --username alice --password secret   # Create user
david user role alice admin                              # Change role (admin)
david user apikey create --name "Desktop"              # Create API key
david user apikey revoke {key_id}                      # Revoke API key

# Admin
david admin stats                                        # Show statistics
david admin audit --since 20260301                       # Show audit log
david admin ratelimit list                               # List rate limits
```

## Success Criteria
- [ ] Import command works (single file, multiple files, dry-run)
- [ ] Export command works (calendar, event, date range, all)
- [ ] Calendar CRUD commands work
- [ ] Calendar share commands work
- [ ] Event list/show/delete commands work
- [ ] User CRUD commands work (admin)
- [ ] API key management commands work
- [ ] Admin commands work (stats, audit, ratelimit)
- [ ] Setup command works
- [ ] CLI integrates with config system
- [ ] 80%+ test coverage
- [ ] All tests pass

## Notes
- CLI should work both locally (direct DB access) and remotely (API client)
- Provide helpful error messages with suggestions
- Use interactive prompts for sensitive operations (delete, password)
- Support both short and long flag formats
- Add help text to all commands and flags

## Next Phase
Move to **Phase 6: WebDAV Integration** when all items above are complete.
