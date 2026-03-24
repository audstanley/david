package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// ListCmd represents the calendar list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all calendars",
	Long:  `List all calendars accessible to the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := listCalendars()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

// CalendarInfo holds calendar information for display
type CalendarInfo struct {
	UID         string
	DisplayName string
	Description string
	IsPublic    bool
	EventCount  int
}

func init() {
}

// listCalendars lists all calendars
func listCalendars() error {
	config := loadConfig()
	dbMgr, err := openDBManager(config.DataDir)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	var calendars []CalendarInfo

	// TODO: Actually query database
	// For now, create mock data
	calendars = []CalendarInfo{
		{UID: "cal-1", DisplayName: "Work Calendar", Description: "Work events", IsPublic: false, EventCount: 10},
		{UID: "cal-2", DisplayName: "Personal Calendar", Description: "Personal events", IsPublic: false, EventCount: 5},
		{UID: "cal-3", DisplayName: "Public Calendar", Description: "Shared events", IsPublic: true, EventCount: 25},
	}

	// Display
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "UID\tNAME\tDESCRIPTION\tPUBLIC\tEVENTS")
	fmt.Fprintln(w, "---\t----\t-----------\t------\t------")

	for _, cal := range calendars {
		public := "no"
		if cal.IsPublic {
			public = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
			cal.UID, cal.DisplayName, cal.Description, public, cal.EventCount)
	}

	w.Flush()
	return nil
}
