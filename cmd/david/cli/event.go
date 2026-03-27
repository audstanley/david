package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/audstanley/david/app/storage"
	"github.com/audstanley/david/app/storage/events"
	"github.com/spf13/cobra"
)

var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "Manage events",
	Long:  `Manage calendar events.`,
}

var EventListCmd = &cobra.Command{
	Use:   "list <calendar-uid>",
	Short: "List events in a calendar",
	Long:  `List all events in a specified calendar.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := listEvents(args[0])
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

var EventShowCmd = &cobra.Command{
	Use:   "show <event-uid>",
	Short: "Show a specific event",
	Long:  `Show details of a specific event by UID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := showEvent(args[0])
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

var EventDeleteCmd = &cobra.Command{
	Use:   "delete <event-uid>",
	Short: "Delete an event",
	Long:  `Delete a specific event by UID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm("Are you sure you want to delete this event?") {
			fmt.Println("Operation cancelled")
			return
		}
		err := deleteEvent(args[0])
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
}

func listEvents(calendarUID string) error {
	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	eventsDB, err := dbMgr.Get(storage.DBEvents)
	if err != nil {
		return fmt.Errorf("failed to get events database: %w", err)
	}

	var eventsList []*storage.Event
	err = eventsDB.Iterate([]byte("event:"), func(key, value []byte) error {
		var e storage.Event
		if err := storage.Decode(value, &e); err != nil {
			return err
		}
		if e.CalendarUID == calendarUID {
			eventsList = append(eventsList, &e)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}

	if len(eventsList) == 0 {
		fmt.Println("No events found in this calendar")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "UID\tSUMMARY\tSTART\tEND")
	fmt.Fprintln(w, "---\t-------\t-----\t---")

	for _, e := range eventsList {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.UID, e.Summary, e.DTStart.Format("2006-01-02 15:04"),
			e.DTEnd.Format("2006-01-02 15:04"))
	}

	w.Flush()
	return nil
}

func showEvent(eventUID string) error {
	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	eventsDB, err := dbMgr.Get(storage.DBEvents)
	if err != nil {
		return fmt.Errorf("failed to get events database: %w", err)
	}

	eventStore := events.New(eventsDB)
	event, err := eventStore.GetByUID(eventUID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	fmt.Printf("Event UID: %s\n", event.UID)
	fmt.Printf("Summary: %s\n", event.Summary)
	fmt.Printf("Description: %s\n", event.Description)
	fmt.Printf("Location: %s\n", event.Location)
	fmt.Printf("Start: %s\n", event.DTStart.Format("2006-01-02 15:04 MST"))
	fmt.Printf("End: %s\n", event.DTEnd.Format("2006-01-02 15:04 MST"))
	fmt.Printf("Calendar: %s\n", event.CalendarUID)
	fmt.Printf("Created: %s\n", event.Created.Format("2006-01-02 15:04 MST"))
	fmt.Printf("Last Modified: %s\n", event.LastModified.Format("2006-01-02 15:04 MST"))

	return nil
}

func deleteEvent(eventUID string) error {
	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	eventsDB, err := dbMgr.Get(storage.DBEvents)
	if err != nil {
		return fmt.Errorf("failed to get events database: %w", err)
	}

	eventStore := events.New(eventsDB)

	if err := eventStore.Delete(eventUID); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	fmt.Printf("Event '%s' deleted successfully\n", eventUID)
	return nil
}
