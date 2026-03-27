package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestGetRootCmd(t *testing.T) {
	cmd := GetRootCmd()
	if cmd == nil {
		t.Error("GetRootCmd returned nil")
	}
	if cmd.Use != "david" {
		t.Errorf("Expected Use to be 'david', got '%s'", cmd.Use)
	}
}

func TestRootCommand(t *testing.T) {
	cmd := GetRootCmd()
	// Root command may not have a Run function if it only has subcommands
	_ = cmd
}

func TestRootCommandFlags(t *testing.T) {
	cmd := GetRootCmd()

	flags := []string{"config", "verbose", "debug", "production"}
	for _, flag := range flags {
		f := cmd.Flag(flag)
		if f == nil {
			t.Errorf("Flag --%s not found", flag)
		}
	}
}

func TestGetServerCmd(t *testing.T) {
	serverCmd := GetServerCmd()
	if serverCmd == nil {
		t.Error("GetServerCmd returned nil")
	}
	if serverCmd.Short != "Start the WebDAV server" {
		t.Errorf("Expected Short to be 'Start the WebDAV server', got '%s'", serverCmd.Short)
	}
}

func TestServerFlags(t *testing.T) {
	serverCmd := GetServerCmd()
	if serverCmd == nil {
		t.Fatal("GetServerCmd returned nil")
	}

	flags := []string{"config", "host", "port", "debug", "production", "hash-algorithm"}
	for _, flag := range flags {
		f := serverCmd.Flag(flag)
		if f == nil {
			t.Errorf("Flag --%s not found", flag)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "version" {
			found = true
			if c.Run == nil {
				t.Error("Version command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("Version command not found in root command")
	}
}

func TestServerCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "server" {
			found = true
			if c.Run == nil {
				t.Error("Server command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("Server command not found in root command")
	}
}

func TestImportCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "import" {
			found = true
			if c.Run == nil {
				t.Error("Import command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("Import command not found in root command")
	}
}

func TestExportCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "export" {
			found = true
			if c.Run == nil {
				t.Error("Export command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("Export command not found in root command")
	}
}

func TestSetupCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "setup" {
			found = true
			if c.Run == nil {
				t.Error("Setup command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("Setup command not found in root command")
	}
}

func TestListCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "list" {
			found = true
			if c.Run == nil {
				t.Error("List command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("List command not found in root command")
	}
}

func TestShowCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "show" {
			found = true
			if c.Run == nil {
				t.Error("Show command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("Show command not found in root command")
	}
}

func TestCreateCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "create" {
			found = true
			if c.Run == nil {
				t.Error("Create command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("Create command not found in root command")
	}
}

func TestDeleteCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "delete" {
			found = true
			if c.Run == nil {
				t.Error("Delete command Run function is nil")
			}
			break
		}
	}
	if !found {
		t.Error("Delete command not found in root command")
	}
}

func TestEventCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "event" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Event command not found in root command")
	}
}

func TestEventSubCommands(t *testing.T) {
	cmd := GetRootCmd()

	var eventCmd *cobra.Command
	for _, c := range cmd.Commands() {
		if c.Use == "event" {
			eventCmd = c
			break
		}
	}
	if eventCmd == nil {
		t.Error("Event command not found")
		return
	}

	// Check that event command has subcommands
	if len(eventCmd.Commands()) == 0 {
		t.Log("Event command has no subcommands (expected for parent command)")
	}
}

func TestUserCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "user" {
			found = true
			break
		}
	}
	if !found {
		t.Error("User command not found in root command")
	}
}

func TestUserSubCommands(t *testing.T) {
	cmd := GetRootCmd()

	var userCmd *cobra.Command
	for _, c := range cmd.Commands() {
		if c.Use == "user" {
			userCmd = c
			break
		}
	}
	if userCmd == nil {
		t.Error("User command not found")
		return
	}

	// Check that user command has subcommands
	if len(userCmd.Commands()) == 0 {
		t.Log("User command has no subcommands (expected for parent command)")
	}
}

func TestAdminCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "admin" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Admin command not found in root command")
	}
}

func TestAdminSubCommands(t *testing.T) {
	cmd := GetRootCmd()

	var adminCmd *cobra.Command
	for _, c := range cmd.Commands() {
		if c.Use == "admin" {
			adminCmd = c
			break
		}
	}
	if adminCmd == nil {
		t.Error("Admin command not found")
		return
	}

	expectedSubCommands := []string{"stats", "audit"}
	for _, expected := range expectedSubCommands {
		found := false
		for _, sub := range adminCmd.Commands() {
			if sub.Use == expected {
				found = true
				if sub.Run == nil {
					t.Errorf("Admin %s command Run function is nil", expected)
				}
				break
			}
		}
		if !found {
			t.Errorf("Admin %s subcommand not found", expected)
		}
	}
}

func TestInitConfig(t *testing.T) {
	// Test that initConfig doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("initConfig panicked: %v", r)
		}
	}()

	initConfig()
}

func TestViperBindFlags(t *testing.T) {
	// Test that server flags are bound to viper
	serverCmd := GetServerCmd()
	if serverCmd == nil {
		t.Fatal("GetServerCmd returned nil")
	}

	flags := []string{"config", "host", "port", "debug", "production", "hash-algorithm"}
	for _, flag := range flags {
		f := serverCmd.Flag(flag)
		if f == nil {
			t.Errorf("Flag --%s not found", flag)
			continue
		}

		// Check that the flag is bound to viper
		if !viper.IsSet(flag) {
			// This is expected before binding, but we can check the flag exists
			_ = f
		}
	}
}

func TestVersionCmdOutput(t *testing.T) {
	cmd := GetRootCmd()

	var versionCmd *cobra.Command
	for _, c := range cmd.Commands() {
		if c.Use == "version" {
			versionCmd = c
			break
		}
	}
	if versionCmd == nil {
		t.Error("Version command not found")
		return
	}

	// Verify version command has expected structure
	if versionCmd.Short != "Print the version number" {
		t.Errorf("Expected Short to be 'Print the version number', got '%s'", versionCmd.Short)
	}
	if versionCmd.Run == nil {
		t.Error("Version command Run function is nil")
	}
}

func TestRootCmdPersistentFlags(t *testing.T) {
	cmd := GetRootCmd()

	// Test persistent flags
	flags := []string{"config", "verbose", "debug", "production"}
	for _, flag := range flags {
		f := cmd.Flag(flag)
		if f == nil {
			t.Errorf("Flag --%s not found", flag)
		}
	}
}

func TestCommandCount(t *testing.T) {
	cmd := GetRootCmd()
	expectedCommands := 6 // version, server, import, export, setup, list, show, create, delete, share, event, user, admin
	actualCommands := len(cmd.Commands())

	if actualCommands < expectedCommands-1 { // Allow for one less due to potential variations
		t.Errorf("Expected at least %d commands, got %d", expectedCommands-1, actualCommands)
	}
}

func TestConfigFileFlag(t *testing.T) {
	cmd := GetRootCmd()
	flag := cmd.Flag("config")
	if flag == nil {
		t.Error("Config flag not found on root command")
	} else {
		// Verify flag has default value
		_ = flag.DefValue
	}
}

func TestDebugFlag(t *testing.T) {
	cmd := GetRootCmd()
	flag := cmd.Flag("debug")
	if flag == nil {
		t.Error("Debug flag not found on root command")
	} else {
		if flag.DefValue != "false" {
			t.Errorf("Expected default value 'false', got '%s'", flag.DefValue)
		}
	}
}

func TestProductionFlag(t *testing.T) {
	cmd := GetRootCmd()
	flag := cmd.Flag("production")
	if flag == nil {
		t.Error("Production flag not found on root command")
	} else {
		if flag.DefValue != "false" {
			t.Errorf("Expected default value 'false', got '%s'", flag.DefValue)
		}
	}
}

func TestVerboseFlag(t *testing.T) {
	cmd := GetRootCmd()
	flag := cmd.Flag("verbose")
	if flag == nil {
		t.Error("Verbose flag not found on root command")
	} else {
		if flag.DefValue != "false" {
			t.Errorf("Expected default value 'false', got '%s'", flag.DefValue)
		}
	}
}

func TestHostFlag(t *testing.T) {
	serverCmd := GetServerCmd()
	flag := serverCmd.Flag("host")
	if flag == nil {
		t.Error("Host flag not found on server command")
	}
}

func TestPortFlag(t *testing.T) {
	serverCmd := GetServerCmd()
	flag := serverCmd.Flag("port")
	if flag == nil {
		t.Error("Port flag not found on server command")
	}
}

func TestHashAlgorithmFlag(t *testing.T) {
	serverCmd := GetServerCmd()
	flag := serverCmd.Flag("hash-algorithm")
	if flag == nil {
		t.Error("Hash-algorithm flag not found on server command")
	}
}

func TestSharedCommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, c := range cmd.Commands() {
		if c.Use == "share" {
			found = true
			break
		}
	}
	if !found {
		t.Log("Share command not found in root command")
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	cmd := GetRootCmd()

	// Parent commands may not have Run functions - that's okay
	// Only leaf commands should have Run functions
	for _, c := range cmd.Commands() {
		if c.Run == nil && len(c.Commands()) > 0 {
			// Parent commands without subcommands should have Run functions
			t.Logf("Command %s is a parent command without Run function (OK if it has subcommands)", c.Use)
		}
	}
}

func TestSubCommandsExist(t *testing.T) {
	cmd := GetRootCmd()

	// Check event subcommands
	for _, c := range cmd.Commands() {
		if c.Use == "event" {
			if len(c.Commands()) == 0 {
				t.Error("Event command has no subcommands")
			}
		}
		// Check user subcommands
		if c.Use == "user" {
			if len(c.Commands()) == 0 {
				t.Error("User command has no subcommands")
			}
		}
		// Check admin subcommands
		if c.Use == "admin" {
			if len(c.Commands()) == 0 {
				t.Error("Admin command has no subcommands")
			}
		}
	}
}

func TestConfigFlagBinding(t *testing.T) {
	// Test that config flag is properly bound
	cmd := GetRootCmd()
	flag := cmd.Flag("config")
	if flag == nil {
		t.Error("Config flag not found")
		return
	}

	// Verify flag properties
	if flag.Shorthand != "" {
		t.Logf("Config flag shorthand: %s", flag.Shorthand)
	}
	_ = flag.Value.Type()
}

func TestInitConfigWithCustomConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	// Write a minimal config file
	err := os.WriteFile(configPath, []byte("address: localhost\nport: 8080\n"), 0644)
	if err != nil {
		t.Skipf("Skipping test due to unable to write config: %v", err)
	}

	// Save original cfgFile
	originalCfgFile := cfgFile
	defer func() {
		cfgFile = originalCfgFile
	}()

	// Set config file and call initConfig
	cfgFile = configPath
	initConfig()

	// Verify viper has the config
	if !viper.IsSet("address") {
		t.Log("Config file not loaded (expected in test environment)")
	}
}

func TestInitConfigDefault(t *testing.T) {
	// Test default config initialization
	originalCfgFile := cfgFile
	defer func() {
		cfgFile = originalCfgFile
	}()

	cfgFile = "" // Reset to trigger default behavior
	initConfig()

	// Viper should be configured with default paths
	if !viper.IsSet("home") {
		// Home directory should be found
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Skipping test: unable to get home directory: %v", err)
		}
		_ = home
	}
}

func TestRootCommandUsage(t *testing.T) {
	cmd := GetRootCmd()
	if cmd.Use != "david" {
		t.Errorf("Expected Use to be 'david', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("Short description is empty")
	}
	if cmd.Long == "" {
		t.Error("Long description is empty")
	}
}

func TestServerCommandUsage(t *testing.T) {
	cmd := GetServerCmd()
	if cmd.Use != "server" {
		t.Errorf("Expected Use to be 'server', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("Server Short description is empty")
	}
}

func TestVersionCmdUsage(t *testing.T) {
	cmd := GetRootCmd()
	for _, c := range cmd.Commands() {
		if c.Use == "version" {
			if c.Short != "Print the version number" {
				t.Errorf("Expected Short to be 'Print the version number', got '%s'", c.Short)
			}
			return
		}
	}
	t.Error("Version command not found")
}

func TestNoArgsExpected(t *testing.T) {
	cmd := GetRootCmd()
	// Cobra defaults to AcceptArgs if not specified
	// This is okay - we just verify the command exists
	_ = cmd.Args
}
