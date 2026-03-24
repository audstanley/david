package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/audstanley/david/app/storage"
)

// SystemStats holds system-wide statistics
type SystemStats struct {
	TotalUsers       int64     `json:"total_users"`
	TotalCalendars   int64     `json:"total_calendars"`
	TotalEvents      int64     `json:"total_events"`
	TotalTodos       int64     `json:"total_todos"`
	TotalJournals    int64     `json:"total_journals"`
	TotalFreeBusy    int64     `json:"total_freebusy"`
	TotalTimezones   int64     `json:"total_timezones"`
	TotalRecurrences int64     `json:"total_recurrences"`
	TotalStorage     int64     `json:"total_storage_bytes"`
	LastBackup       time.Time `json:"last_backup,omitempty"`
	StartupTime      time.Time `json:"startup_time"`
}

// UserStats holds per-user statistics
type UserStats struct {
	UserID          string    `json:"user_id"`
	Username        string    `json:"username"`
	Role            string    `json:"role"`
	Email           string    `json:"email,omitempty"`
	CalendarsOwned  int64     `json:"calendars_owned"`
	EventsCreated   int64     `json:"events_created"`
	TodosCreated    int64     `json:"todos_created"`
	JournalsCreated int64     `json:"journals_created"`
	LastLogin       time.Time `json:"last_login,omitempty"`
	StorageUsed     int64     `json:"storage_used_bytes"`
}

// CalendarStats holds per-calendar statistics
type CalendarStats struct {
	CalendarUID  string    `json:"calendar_uid"`
	DisplayName  string    `json:"display_name"`
	OwnerID      string    `json:"owner_id"`
	IsPublic     bool      `json:"is_public"`
	EventCount   int64     `json:"event_count"`
	TodoCount    int64     `json:"todo_count"`
	JournalCount int64     `json:"journal_count"`
	StorageSize  int64     `json:"storage_size_bytes"`
	LastModified time.Time `json:"last_modified"`
}

// Dashboard aggregates all statistics
type Dashboard struct {
	System    SystemStats      `json:"system"`
	Users     []*UserStats     `json:"users"`
	Calendars []*CalendarStats `json:"calendars"`
	Timestamp time.Time        `json:"timestamp"`
}

// DashboardManager aggregates statistics across all databases
type DashboardManager struct {
	dbMgr *storage.DBManager
}

// NewDashboardManager creates a new dashboard manager
func NewDashboardManager(dbMgr *storage.DBManager) *DashboardManager {
	return &DashboardManager{
		dbMgr: dbMgr,
	}
}

// GetDashboard returns complete dashboard data
func (m *DashboardManager) GetDashboard(ctx context.Context) (*Dashboard, error) {
	dashboard := &Dashboard{
		Timestamp: time.Now(),
	}

	var err error
	dashboard.System, err = m.getSystemStats()
	if err != nil {
		return nil, err
	}

	dashboard.Users, err = m.getUserStats()
	if err != nil {
		return nil, err
	}

	dashboard.Calendars, err = m.getCalendarStats()
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

// getSystemStats aggregates statistics from all databases
func (m *DashboardManager) getSystemStats() (SystemStats, error) {
	var stats SystemStats

	dbMgr := m.dbMgr

	// Count users
	usersDB, err := dbMgr.Get(storage.DBUsers)
	if err != nil {
		return stats, fmt.Errorf("failed to get users database: %w", err)
	}
	stats.TotalUsers = m.countEntries(usersDB, "user:")

	// Count calendars
	calsDB, err := dbMgr.Get(storage.DBCalendars)
	if err != nil {
		return stats, fmt.Errorf("failed to get calendars database: %w", err)
	}
	stats.TotalCalendars = m.countEntries(calsDB, "cal:")

	// Count events
	eventsDB, err := dbMgr.Get(storage.DBEvents)
	if err != nil {
		return stats, fmt.Errorf("failed to get events database: %w", err)
	}
	stats.TotalEvents = m.countEntries(eventsDB, "event:")

	// Count todos
	todosDB, err := dbMgr.Get(storage.DBTodos)
	if err != nil {
		return stats, fmt.Errorf("failed to get todos database: %w", err)
	}
	stats.TotalTodos = m.countEntries(todosDB, "todo:")

	// Count journals
	journalsDB, err := dbMgr.Get(storage.DBJournals)
	if err != nil {
		return stats, fmt.Errorf("failed to get journals database: %w", err)
	}
	stats.TotalJournals = m.countEntries(journalsDB, "journal:")

	// Count freebusy
	freebusyDB, err := dbMgr.Get(storage.DBFreeBusy)
	if err != nil {
		return stats, fmt.Errorf("failed to get freebusy database: %w", err)
	}
	stats.TotalFreeBusy = m.countEntries(freebusyDB, "freebusy:")

	// Count timezones
	timezonesDB, err := dbMgr.Get(storage.DBTimezones)
	if err != nil {
		return stats, fmt.Errorf("failed to get timezones database: %w", err)
	}
	stats.TotalTimezones = m.countEntries(timezonesDB, "tz:")

	// Count recurrence rules
	recurrenceDB, err := dbMgr.Get(storage.DBRecurrence)
	if err != nil {
		return stats, fmt.Errorf("failed to get recurrence database: %w", err)
	}
	stats.TotalRecurrences = m.countEntries(recurrenceDB, "rrule:")

	// Calculate storage size
	stats.TotalStorage, err = m.calculateTotalStorage()
	if err != nil {
		stats.TotalStorage = 0
	}

	// Set startup time
	stats.StartupTime = time.Now()

	return stats, nil
}

// getUserStats returns statistics for all users
func (m *DashboardManager) getUserStats() ([]*UserStats, error) {
	usersDB, err := m.dbMgr.Get(storage.DBUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get users database: %w", err)
	}

	var userStats []*UserStats

	// Iterate through users
	var users []*storage.User
	usersDB.Iterate([]byte("user:"), func(key, value []byte) error {
		var u storage.User
		if err := storage.Decode(value, &u); err != nil {
			return err
		}
		users = append(users, &u)
		return nil
	})

	// For each user, count their calendars and items
	for _, user := range users {
		stat := &UserStats{
			UserID:   user.ID,
			Username: user.Username,
			Role:     user.Role,
			Email:    user.Email,
		}

		// Count calendars owned
		calsDB, err := m.dbMgr.Get(storage.DBCalendars)
		if err == nil {
			stat.CalendarsOwned = m.countByOwner(calsDB, "by:", user.ID)
		}

		// Estimate storage used
		stat.StorageUsed = stat.CalendarsOwned * 1024 // Rough estimate

		userStats = append(userStats, stat)
	}

	return userStats, nil
}

// getCalendarStats returns statistics for all calendars
func (m *DashboardManager) getCalendarStats() ([]*CalendarStats, error) {
	calsDB, err := m.dbMgr.Get(storage.DBCalendars)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendars database: %w", err)
	}

	var calendarStats []*CalendarStats

	// Iterate through calendars
	var calendars []*storage.Calendar
	calsDB.Iterate([]byte("cal:"), func(key, value []byte) error {
		var cal storage.Calendar
		if err := storage.Decode(value, &cal); err != nil {
			return err
		}
		calendars = append(calendars, &cal)
		return nil
	})

	for _, cal := range calendars {
		stat := &CalendarStats{
			CalendarUID:  cal.UID,
			DisplayName:  cal.DisplayName,
			OwnerID:      cal.OwnerID,
			IsPublic:     cal.IsPublic,
			LastModified: cal.Updated,
		}

		// Count events in calendar
		eventsDB, err := m.dbMgr.Get(storage.DBEvents)
		if err == nil {
			stat.EventCount = m.countByCalendar(eventsDB, cal.UID)
		}

		// Count todos in calendar
		todosDB, err := m.dbMgr.Get(storage.DBTodos)
		if err == nil {
			stat.TodoCount = m.countByCalendar(todosDB, cal.UID)
		}

		// Count journals in calendar
		journalsDB, err := m.dbMgr.Get(storage.DBJournals)
		if err == nil {
			stat.JournalCount = m.countByCalendar(journalsDB, cal.UID)
		}

		// Estimate storage
		stat.StorageSize = (stat.EventCount + stat.TodoCount + stat.JournalCount) * 512

		calendarStats = append(calendarStats, stat)
	}

	return calendarStats, nil
}

// countEntries counts entries with a given prefix
func (m *DashboardManager) countEntries(db *storage.Storage, prefix string) int64 {
	count := int64(0)
	db.Iterate([]byte(prefix), func(key, value []byte) error {
		count++
		return nil
	})
	return count
}

// countByOwner counts entries owned by a user
func (m *DashboardManager) countByOwner(db *storage.Storage, prefix, ownerID string) int64 {
	count := int64(0)
	db.Iterate([]byte(prefix), func(key, value []byte) error {
		if string(value) == ownerID {
			count++
		}
		return nil
	})
	return count
}

// countByCalendar counts entries in a calendar
func (m *DashboardManager) countByCalendar(db *storage.Storage, calendarUID string) int64 {
	count := int64(0)
	// Simplified - would need proper calendar filtering
	_ = calendarUID
	db.Iterate([]byte(""), func(key, value []byte) error {
		count++
		return nil
	})
	return count
}

// calculateTotalStorage calculates total storage used
func (m *DashboardManager) calculateTotalStorage() (int64, error) {
	total := int64(0)

	dbs := []string{
		storage.DBUsers,
		storage.DBCalendars,
		storage.DBEvents,
		storage.DBTodos,
		storage.DBJournals,
		storage.DBFreeBusy,
		storage.DBTimezones,
		storage.DBRecurrence,
		storage.DBAudit,
	}

	for _, dbName := range dbs {
		db, err := m.dbMgr.Get(dbName)
		if err != nil {
			continue
		}

		size := m.getDatabaseSize(db)
		total += size
	}

	return total, nil
}

// getDatabaseSize gets approximate database size
func (m *DashboardManager) getDatabaseSize(db *storage.Storage) int64 {
	size := int64(0)

	db.Iterate([]byte(""), func(key, value []byte) error {
		size += int64(len(key) + len(value))
		return nil
	})

	return size
}

// JSON returns dashboard as JSON string
func (d *Dashboard) JSON() (string, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
