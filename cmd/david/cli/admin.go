package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/audstanley/david/app/storage"
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin commands",
	Long:  `Administrative commands for david.`,
}

var AdminStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show system statistics",
	Long:  `Show system statistics including user counts, calendar counts, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := showStats()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

var AdminAuditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Show audit logs",
	Long:  `Show recent audit log entries.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := showAudit()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
}

func showStats() error {
	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	fmt.Println("=== DAVID System Statistics ===")
	fmt.Println()

	// Users count
	usersDB, err := dbMgr.Get(storage.DBUsers)
	if err != nil {
		return fmt.Errorf("failed to get users database: %w", err)
	}
	userCount := 0
	usersDB.Iterate([]byte("user:"), func(key, value []byte) error {
		userCount++
		return nil
	})
	fmt.Printf("Users: %d\n", userCount)

	// Calendars count
	calsDB, err := dbMgr.Get(storage.DBCalendars)
	if err != nil {
		return fmt.Errorf("failed to get calendars database: %w", err)
	}
	calCount := 0
	calsDB.Iterate([]byte("cal:"), func(key, value []byte) error {
		calCount++
		return nil
	})
	fmt.Printf("Calendars: %d\n", calCount)

	// Events count
	eventsDB, err := dbMgr.Get(storage.DBEvents)
	if err != nil {
		return fmt.Errorf("failed to get events database: %w", err)
	}
	eventCount := 0
	eventsDB.Iterate([]byte("event:"), func(key, value []byte) error {
		eventCount++
		return nil
	})
	fmt.Printf("Events: %d\n", eventCount)

	// Todos count
	todosDB, err := dbMgr.Get(storage.DBTodos)
	if err != nil {
		return fmt.Errorf("failed to get todos database: %w", err)
	}
	todoCount := 0
	todosDB.Iterate([]byte("todo:"), func(key, value []byte) error {
		todoCount++
		return nil
	})
	fmt.Printf("Todos: %d\n", todoCount)

	// Journals count
	journalsDB, err := dbMgr.Get(storage.DBJournals)
	if err != nil {
		return fmt.Errorf("failed to get journals database: %w", err)
	}
	journalCount := 0
	journalsDB.Iterate([]byte("journal:"), func(key, value []byte) error {
		journalCount++
		return nil
	})
	fmt.Printf("Journals: %d\n", journalCount)

	// Recurrence rules count
	recurrenceDB, err := dbMgr.Get(storage.DBRecurrence)
	if err != nil {
		return fmt.Errorf("failed to get recurrence database: %w", err)
	}
	recurrenceCount := 0
	recurrenceDB.Iterate([]byte("rrule:"), func(key, value []byte) error {
		recurrenceCount++
		return nil
	})
	fmt.Printf("Recurrence Rules: %d\n", recurrenceCount)

	// Timezones count
	timezonesDB, err := dbMgr.Get(storage.DBTimezones)
	if err != nil {
		return fmt.Errorf("failed to get timezones database: %w", err)
	}
	tzCount := 0
	timezonesDB.Iterate([]byte("tz:"), func(key, value []byte) error {
		tzCount++
		return nil
	})
	fmt.Printf("Timezones: %d\n", tzCount)

	return nil
}

func showAudit() error {
	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	auditDB, err := dbMgr.Get(storage.DBAudit)
	if err != nil {
		return fmt.Errorf("failed to get audit database: %w", err)
	}

	var entries []string
	auditDB.Iterate([]byte("audit:"), func(key, value []byte) error {
		entries = append(entries, string(value))
		return nil
	})

	if len(entries) == 0 {
		fmt.Println("No audit log entries found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tUSER\tACTION\tDETAILS")
	fmt.Fprintln(w, "---------\t----\t------\t-------")

	for i := len(entries) - 1; i >= 0 && i > len(entries)-20; i-- {
		fmt.Fprintln(w, entries[i])
	}

	w.Flush()
	return nil
}
