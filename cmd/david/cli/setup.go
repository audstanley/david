package cli

import (
	"fmt"
	"os"

	"github.com/audstanley/david/app/storage"
	"github.com/spf13/cobra"
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize david",
	Long: `Initialize david by creating necessary directories and databases.

This command sets up the data directory structure and creates all required databases.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := setupDavid()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
}

func setupDavid() error {
	dataDir := "./data/david"

	// Create data directory
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	fmt.Printf("Created data directory: %s\n", dataDir)

	// Create database manager
	dbMgr, err := storage.NewDBManager(dataDir)
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}
	defer dbMgr.CloseAll()

	// Create all databases
	if err := dbMgr.CreateAll(); err != nil {
		return fmt.Errorf("failed to create databases: %w", err)
	}

	fmt.Println("Created all databases:")
	fmt.Println("  - Users database")
	fmt.Println("  - Calendars database")
	fmt.Println("  - Events database")
	fmt.Println("  - Todos database")
	fmt.Println("  - Journals database")
	fmt.Println("  - Freebusy database")
	fmt.Println("  - Timezones database")
	fmt.Println("  - Recurrence database")
	fmt.Println("  - Audit database")

	fmt.Println()
	fmt.Println("Setup complete!")
	fmt.Println()
	fmt.Println("You can now create your first user:")
	fmt.Println("  david user create <username>")

	return nil
}
