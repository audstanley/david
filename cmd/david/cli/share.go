package cli

import (
	"fmt"
	"os"

	"github.com/audstanley/david/app/storage"
	"github.com/audstanley/david/app/storage/calendars"
	"github.com/audstanley/david/app/storage/users"
	"github.com/spf13/cobra"
)

var ShareCmd = &cobra.Command{
	Use:   "share <calendar-uid> <user-id>",
	Short: "Share a calendar with a user",
	Long: `Share a calendar with another user by their user ID.

Example:
  david calendar share cal-123 user-456`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := shareCalendar(args[0], args[1])
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func shareCalendar(calendarUID, userID string) error {
	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	calsDB, err := dbMgr.Get(storage.DBCalendars)
	if err != nil {
		return fmt.Errorf("failed to get calendars database: %w", err)
	}

	usersDB, err := dbMgr.Get(storage.DBUsers)
	if err != nil {
		return fmt.Errorf("failed to get users database: %w", err)
	}

	clients := &Clients{
		CalendarStore: calendars.New(calsDB),
		UserStore:     users.New(usersDB),
	}

	calendar, err := clients.CalendarStore.GetByUID(calendarUID)
	if err != nil {
		return fmt.Errorf("failed to get calendar: %w", err)
	}

	user, err := clients.UserStore.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if calendar.OwnerID == userID {
		return fmt.Errorf("cannot share calendar with its owner")
	}

	if err := clients.CalendarStore.Update(calendarUID, map[string]interface{}{
		"isPublic": true,
	}); err != nil {
		return fmt.Errorf("failed to update calendar: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Calendar '%s' is now public and accessible by user '%s'\n", calendar.DisplayName, user.Username)
	return nil
}
