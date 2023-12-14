package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestParseConfig(t *testing.T) {
	// Reset Viper to ensure clean state across tests
	viper.Reset()

	// Create a temporary directory to store test config files
	tmpDir := filepath.Join(os.TempDir(), "david__"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(tmpDir, 0700)
	defer os.RemoveAll(tmpDir) // Automatically clean up temp directory after tests

	// Define test cases with expected configurations
	tests := []struct {
		name string  // Test case name
		want *Config // Expected configuration after parsing (created by cfg function)
	}{
		{"default", cfg(t, tmpDir)}, // Test default config loaded from temp directory
	}
	// Loop through each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the configuration with an empty path (use config in temp dir)
			got := ParseConfig("")
			// Compare the parsed config with the expected config
			if !reflect.DeepEqual(got, tt.want) {
				// Marshal both configs to JSON for easier comparison in error message
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.want)
				t.Errorf("ParseConfig(\"\") = %s, want %s", gotJSON, wantJSON)
			}
		})
	}
}

func cfg(t *testing.T, tmpDir string) *Config {
	// **1. Test Configuration Setup**

	// **a. Set the configuration format to YAML:**
	viper.SetConfigType("yaml")
	// **b. Define the test YAML configuration as a byte slice:**
	// This variable holds the test configuration data in a compact format.
	var yamlCfg = []byte(`
address: 1.2.3.4
port: 42
prefix: /oh-de-lally
tls:
  keyFile: ` + tmpDir + `/robin.pem
  certFile: ` + tmpDir + `/tuck.pem
dir: /sherwood/forest
realm: uk
users:
  lj:
    password: 123
    subdir: /littlejohn
    permissions: crud
  srf:
    password: 234
    subdir: /sheriff
    permissions: crud
log:
  error: true
`)
	// **2. Temporary File Creation**

	// **a. Create temporary files for TLS key and certificate:**
	// These files will be populated with empty data to represent the TLS configuration.
	err := os.WriteFile(filepath.Join(tmpDir, "config.yaml"), yamlCfg, 0600)
	if err != nil {
		t.Errorf("error writing test config. error = %v", err)
	}

	// **3. Configuration Reading and Parsing**

	// **a. Read the configuration from the temporary file:**
	err = viper.ReadConfig(bytes.NewBuffer(yamlCfg))
	if err != nil {
		t.Errorf("error reading test config. error = %v", err)
	}
	// **b. Allocate and unmarshal the configuration data:**
	// Create a new Config instance and populate it with the parsed data.
	var cfg = &Config{}
	viper.Unmarshal(&cfg)

	// **4. User Permissions Processing**
	for user := range viper.GetStringMap("Users") {
		permissions := viper.GetString(fmt.Sprintf("Users.%s.permissions", user)) // Access specific user permissions
		cfg.Users[user].Crud = &CrudType{Crud: permissions}                       // Set user's CRUD permissions object
		err := FormatCrud(context.Background(), user, cfg)                        // Further process and validate permissions
		if err != nil {
			log.WithError(err).WithField("user", user).Error("Error parsing crud string from config file") // log error with context
		}
		log.WithFields(logrus.Fields{"user": user,
			"crud": cfg.Users[user].Crud}).Info("Parsed crud string from config file") // Log parsed permissions
	}

	// **5. Config Path and Dummy Files (Optional)**

	// **a. Add the temporary directory to the config path:**
	// This ensures viper prioritizes the test configuration.
	viper.AddConfigPath(tmpDir)

	// add dummy cert and key file
	_, err = os.OpenFile(filepath.Join(tmpDir, "robin.pem"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Errorf("error creating key file. error = %v", err)
		return nil
	}

	_, err = os.OpenFile(filepath.Join(tmpDir, "tuck.pem"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Errorf("error creating cert file. error = %v", err)
		return nil
	}
	viper.AddConfigPath(tmpDir)
	// **6. Return the Config Instance**
	// Return the populated Config instance for further use in the test case.
	return cfg
}
