package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	calcName      string
	calcDesc      string
	calcColor     string
	calcPublic    bool
	calcTimezone  string
)

// CreateCmd represents the calendar create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new calendar",
	Long:  `Create a new calendar with the specified properties.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := createCalendar()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {

	CreateCmd.Flags().StringVarP(&calcName, "name", "n", "", "Calendar display name (required)")
	CreateCmd.Flags().StringVar(&calcDesc, "description", "", "Calendar description")
	CreateCmd.Flags().StringVar(&calcColor, "color", "", "Calendar color (e.g., #FF0000)")
	CreateCmd.Flags().BoolVarP(&calcPublic, "public", "p", false, "Make calendar public")
	CreateCmd.Flags().StringVar(&calcTimezone, "timezone", "UTC", "Calendar timezone")
	
	CreateCmd.MarkFlagRequired("name")
}

// createCalendar creates a new calendar
func createCalendar() error {
	if calcName == "" {
		return fmt.Errorf("--name is required")
	}

	// TODO: Actually create calendar in database
	// For now, just print success
	fmt.Printf("Created calendar: %s\n", calcName)
	fmt.Printf("  UID: cal-%d\n", 1234567890) // Mock UID
	fmt.Printf("  Name: %s\n", calcName)
	if calcDesc != "" {
		fmt.Printf("  Description: %s\n", calcDesc)
	}
	if calcColor != "" {
		fmt.Printf("  Color: %s\n", calcColor)
	}
	fmt.Printf("  Public: %v\n", calcPublic)
	
	return nil
}
