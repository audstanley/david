package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	verbose    bool
	debug      bool
	production bool
	version    = "0.1.0"
	buildTime  = "unknown"
	gitCommit  = "unknown"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "david",
	Short: "david - A CalDAV server",
	Long: `david is a CalDAV server that provides calendar, event, and todo management.

It supports:
- Full iCalendar (RFC 2445) parsing and generation
- Multi-user access with RBAC
- REST API and WebDAV/CalDAV protocol support
- Calendar sharing and collaboration
- Recurring events and timezones`,
}

// Execute adds all child commands to the root command
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.david.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	RootCmd.PersistentFlags().BoolVarP(&production, "production", "P", false, "production mode")

	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(serverCmd)
	RootCmd.AddCommand(ImportCmd)
	RootCmd.AddCommand(ExportCmd)
	RootCmd.AddCommand(SetupCmd)

	// Calendar commands
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(ShowCmd)
	RootCmd.AddCommand(CreateCmd)
	RootCmd.AddCommand(DeleteCmd)
	RootCmd.AddCommand(ShareCmd)

	// Event commands
	RootCmd.AddCommand(EventCmd)
	EventCmd.AddCommand(EventListCmd)
	EventCmd.AddCommand(EventShowCmd)
	EventCmd.AddCommand(EventDeleteCmd)

	// User commands
	RootCmd.AddCommand(UserCmd)
	UserCmd.AddCommand(UserListCmd)
	UserCmd.AddCommand(UserCreateCmd)
	UserCmd.AddCommand(UserDeleteCmd)

	// Admin commands
	RootCmd.AddCommand(AdminCmd)
	AdminCmd.AddCommand(AdminStatsCmd)
	AdminCmd.AddCommand(AdminAuditCmd)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".david" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".david")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
		}
	}
}

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. This is david's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("david version %s (built: %s, commit: %s)\n", version, buildTime, gitCommit)
	},
}
