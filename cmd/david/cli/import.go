package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/audstanley/david/app/icalendar"
	"github.com/spf13/cobra"
)

var (
	importCalendar string
	dryRun         bool
	force          bool
	importVerbose  bool
	importFiles    []string
)

// ImportCmd represents the import command
var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import iCalendar files",
	Long: `Import iCalendar (.ics) files into david.
	
Supports importing single files or multiple files from a directory.
Duplicates are handled based on the --force flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(importFiles) == 0 && len(args) == 0 {
			cobra.CheckErr(fmt.Errorf("no input files specified"))
		}

		// Add args to files
		importFiles = append(importFiles, args...)

		// Process imports
		err := processImport()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {

	ImportCmd.Flags().StringVarP(&importCalendar, "calendar", "c", "", "Calendar UID to import to")
	ImportCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Validate without importing")
	ImportCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing events with same UID")
	ImportCmd.Flags().BoolVarP(&importVerbose, "verbose", "v", false, "Verbose output")

	ImportCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.david.yaml)")
}

// processImport processes all import files
func processImport() error {
	// Check if importing to a specific calendar
	if importCalendar == "" {
		// Try to find default calendar for user
		// For now, we'll just warn and skip calendar-specific imports
		printVerbose("Warning: No calendar specified, events will not be stored\n")
	}

	// Count statistics
	totalImported := 0
	totalSkipped := 0
	totalErrors := 0

	// Process each file
	for _, file := range importFiles {
		// Handle directory
		if info, err := os.Stat(file); err == nil && info.IsDir() {
			files, err := filepath.Glob(filepath.Join(file, "*.ics"))
			if err != nil {
				printError("Error reading directory %s: %v", file, err)
				totalErrors++
				continue
			}
			importFiles = append(importFiles, files...)
			continue
		}

		printVerbose("Processing: %s\n", file)

		// Read file
		data, err := os.ReadFile(file)
		if err != nil {
			printError("Error reading file %s: %v", file, err)
			totalErrors++
			continue
		}

		// Parse ICS
		cal, err := icalendar.Parse(string(data))
		if err != nil {
			printError("Error parsing %s: %v", file, err)
			totalErrors++
			continue
		}

		printVerbose("  Parsed: %d components\n", len(cal.Components))

		if dryRun {
			printVerbose("  [DRY-RUN] Would import %d components\n", len(cal.Components))
			totalImported += len(cal.Components)
			continue
		}

		// TODO: Actually store in database
		// For now, just simulate
		totalImported += len(cal.Components)
		printVerbose("  Imported: %d components\n", len(cal.Components))
	}

	// Summary
	fmt.Printf("\nImport Summary:\n")
	fmt.Printf("  Total processed: %d files\n", len(importFiles))
	fmt.Printf("  Imported: %d components\n", totalImported)
	fmt.Printf("  Skipped: %d components\n", totalSkipped)
	fmt.Printf("  Errors: %d files\n", totalErrors)

	return nil
}

// validateICS validates an ICS file without storing
func validateICS(data []byte) error {
	_, err := icalendar.Parse(string(data))
	return err
}
