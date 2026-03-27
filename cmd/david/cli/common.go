package cli

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/audstanley/david/app/storage"
	"github.com/audstanley/david/app/storage/audit"
	"github.com/audstanley/david/app/storage/calendars"
	"github.com/audstanley/david/app/storage/events"
	"github.com/audstanley/david/app/storage/journals"
	"github.com/audstanley/david/app/storage/recurrence"
	"github.com/audstanley/david/app/storage/timezones"
	"github.com/audstanley/david/app/storage/todos"
	"github.com/audstanley/david/app/storage/users"
)

// Config holds CLI configuration
type Config struct {
	DataDir string
}

// loadConfig loads configuration from viper
func loadConfig() *Config {
	c := &Config{
		DataDir: viper.GetString("data_dir"),
	}

	if c.DataDir == "" {
		c.DataDir = "./data/david"
	}

	return c
}

// openDBManager opens all LevelDB databases
func openDBManager(dataDir string) (*storage.DBManager, error) {
	dbMgr, err := storage.NewDBManager(dataDir)
	if err != nil {
		return nil, err
	}
	if err := dbMgr.CreateAll(); err != nil {
		return nil, err
	}
	return dbMgr, nil
}

// createClients creates storage clients for all databases
func createClients(dbMgr *storage.DBManager) (*Clients, error) {
	usersDB, err := dbMgr.Get(storage.DBUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get users database: %w", err)
	}

	calsDB, err := dbMgr.Get(storage.DBCalendars)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendars database: %w", err)
	}

	eventsDB, err := dbMgr.Get(storage.DBEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to get events database: %w", err)
	}

	todosDB, err := dbMgr.Get(storage.DBTodos)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos database: %w", err)
	}

	journalsDB, err := dbMgr.Get(storage.DBJournals)
	if err != nil {
		return nil, fmt.Errorf("failed to get journals database: %w", err)
	}

	recurrenceDB, err := dbMgr.Get(storage.DBRecurrence)
	if err != nil {
		return nil, fmt.Errorf("failed to get recurrence database: %w", err)
	}

	timezonesDB, err := dbMgr.Get(storage.DBTimezones)
	if err != nil {
		return nil, fmt.Errorf("failed to get timezones database: %w", err)
	}

	auditDB, err := dbMgr.Get(storage.DBAudit)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit database: %w", err)
	}

	return &Clients{
		UserStore:       users.New(usersDB),
		CalendarStore:   calendars.New(calsDB),
		EventStore:      events.New(eventsDB),
		TodoStore:       todos.New(todosDB),
		JournalStore:    journals.New(journalsDB),
		RecurrenceStore: recurrence.New(recurrenceDB),
		TimeZoneStore:   timezones.New(timezonesDB),
		AuditStore:      audit.New(auditDB),
	}, nil
}

// Clients holds all storage clients
type Clients struct {
	UserStore       *users.Store
	CalendarStore   *calendars.Store
	EventStore      *events.Store
	TodoStore       *todos.Store
	JournalStore    *journals.Store
	RecurrenceStore *recurrence.Store
	TimeZoneStore   *timezones.Store
	AuditStore      *audit.Store
}

// printVerbose prints message if verbose mode is enabled
func printVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format, args...)
	}
}

// printError prints error message
func printError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}

// confirm asks for user confirmation
func confirm(prompt string) bool {
	var response string
	fmt.Printf("%s [y/N]: ", prompt)
	fmt.Scan(&response)
	return response == "y" || response == "Y"
}
