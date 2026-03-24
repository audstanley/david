package storage

import (
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// DBManager manages multiple LevelDB databases
type DBManager struct {
	mu        sync.RWMutex
	databases map[string]*Storage
	basePath  string
}

// NewDBManager creates a new database manager
func NewDBManager(basePath string) (*DBManager, error) {
	mgr := &DBManager{
		databases: make(map[string]*Storage),
		basePath:  basePath,
	}
	return mgr, nil
}

// Database names
const (
	DBUsers      = "users"
	DBCalendars  = "calendars"
	DBEvents     = "events"
	DBTodos      = "todos"
	DBJournals   = "journals"
	DBFreeBusy   = "freebusy"
	DBTimezones  = "timezones"
	DBRecurrence = "recurrence"
	DBAudit      = "audit"
)

// AllDatabases returns all database names
func AllDatabases() []string {
	return []string{
		DBUsers,
		DBCalendars,
		DBEvents,
		DBTodos,
		DBJournals,
		DBFreeBusy,
		DBTimezones,
		DBRecurrence,
		DBAudit,
	}
}

// OpenAll opens all databases
func (m *DBManager) OpenAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, name := range AllDatabases() {
		if _, exists := m.databases[name]; exists {
			continue
		}
		if err := m.openDatabase(name); err != nil {
			return fmt.Errorf("failed to open database %s: %w", name, err)
		}
	}
	return nil
}

// CloseAll closes all databases
func (m *DBManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for name, db := range m.databases {
		if err := db.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close database %s: %w", name, err)
		}
		delete(m.databases, name)
	}
	return lastErr
}

// openDatabase opens a single database
func (m *DBManager) openDatabase(name string) error {
	path := fmt.Sprintf("%s/%s", m.basePath, name)
	db, err := NewStorage(path)
	if err != nil {
		return err
	}
	m.databases[name] = db
	return nil
}

// Get returns a database by name
func (m *DBManager) Get(name string) (*Storage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	db, exists := m.databases[name]
	if !exists {
		return nil, &StorageError{Code: "NOT_FOUND", Message: fmt.Sprintf("database %s not found", name)}
	}
	return db, nil
}

// Has checks if a database exists
func (m *DBManager) Has(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.databases[name]
	return exists
}

// CreateAll creates all databases if they don't exist
func (m *DBManager) CreateAll() error {
	for _, name := range AllDatabases() {
		if !m.Has(name) {
			if err := m.openDatabase(name); err != nil {
				return err
			}
		}
	}
	return nil
}

// Path returns the path for a database
func (m *DBManager) Path(name string) string {
	return fmt.Sprintf("%s/%s", m.basePath, name)
}

// Compact compacts all databases
func (m *DBManager) Compact() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lastErr error
	for name, db := range m.databases {
		if err := db.db.CompactRange(util.Range{}); err != nil {
			lastErr = fmt.Errorf("failed to compact database %s: %w", name, err)
		}
	}
	return lastErr
}

// Backup creates a backup of all databases
func (m *DBManager) Backup(backupDir string) error {
	// TODO: Implement backup using leveldb.Export
	return &StorageError{Code: "NOT_IMPLEMENTED", Message: "backup not yet implemented"}
}

// Restore restores databases from a backup
func (m *DBManager) Restore(backupDir string) error {
	// TODO: Implement restore using leveldb.Import
	return &StorageError{Code: "NOT_IMPLEMENTED", Message: "restore not yet implemented"}
}

// Stats returns statistics for all databases
func (m *DBManager) Stats() (map[string]*leveldb.DBStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]*leveldb.DBStats)
	for name, db := range m.databases {
		var dbStats leveldb.DBStats
		if err := db.db.Stats(&dbStats); err != nil {
			return nil, fmt.Errorf("failed to get stats for %s: %w", name, err)
		}
		stats[name] = &dbStats
	}
	return stats, nil
}
