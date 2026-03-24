package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	calcUID string
)

// ShowCmd represents the calendar show command
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show calendar details",
	Long:  `Show detailed information about a calendar.`,
	Run: func(cmd *cobra.Command, args []string) {
		if calcUID == "" && len(args) > 0 {
			calcUID = args[0]
		}
		err := showCalendar()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	ShowCmd.Flags().StringVarP(&calcUID, "uid", "u", "", "Calendar UID")
	ShowCmd.MarkFlagRequired("uid")
}

// showCalendar shows calendar details
func showCalendar() error {
	// TODO: Query database
	fmt.Printf("Calendar: %s\n", calcUID)
	fmt.Printf("  Name: Work Calendar\n")
	fmt.Printf("  Description: Work events\n")
	fmt.Printf("  Color: #3788d8\n")
	fmt.Printf("  Public: no\n")
	fmt.Printf("  Timezone: America/New_York\n")
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nEvents:")
	fmt.Fprintln(w, "UID\tSUMMARY\tSTART\tEND")
	fmt.Fprintln(w, "---\t-------\t-----\t---")
	fmt.Fprintf(w, "event-1\tTeam Meeting\t2026-03-25T14:00\t2026-03-25T15:00\n")
	w.Flush()
	
	return nil
}
