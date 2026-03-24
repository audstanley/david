package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	calcDeleteConfirm bool
)

// DeleteCmd represents the calendar delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a calendar",
	Long:  `Delete a calendar and all its events.`,
	Run: func(cmd *cobra.Command, args []string) {
		if calcUID == "" && len(args) > 0 {
			calcUID = args[0]
		}
		if calcUID == "" {
			cobra.CheckErr(fmt.Errorf("--uid or argument required"))
		}
		
		if !calcDeleteConfirm {
			if !confirm("Are you sure you want to delete calendar " + calcUID + "?") {
				fmt.Println("Cancelled")
				return
			}
		}
		
		err := deleteCalendar()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	DeleteCmd.Flags().StringVarP(&calcUID, "uid", "u", "", "Calendar UID")
	DeleteCmd.Flags().BoolVarP(&calcDeleteConfirm, "confirm", "y", false, "Skip confirmation prompt")
}

// deleteCalendar deletes a calendar
func deleteCalendar() error {
	// TODO: Actually delete from database
	fmt.Printf("Deleted calendar: %s\n", calcUID)
	return nil
}
