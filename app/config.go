package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config represents the configuration of the server application.
type Config struct {
	Address string               `default:"127.0.0.1"`
	Port    string               `default:"8000"`
	Prefix  string               `default:""`
	Dir     string               `default:"/tmp"`
	TLS     *TLS                 `default:"nil"`
	Log     Logging              `default:"{error:true, create:false, read:false, update:false, delete:false}"`
	Realm   string               `default:"david"`
	Users   map[string]*UserInfo `default:"nil"`
	Cors    Cors                 `default:"{origin:*, credentials:false}"`
}

// Logging allows definition for logging each CRUD method.
type Logging struct {
	Production bool `default:"false"`
	Debug      bool `default:"true"`
	Error      bool
	Create     bool
	Read       bool
	Update     bool
	Delete     bool
}

// TLS allows specification of a certificate and private key file.
type TLS struct {
	CertFile string
	KeyFile  string
}

// UserInfo allows storing of a password and user directory.
type UserInfo struct {
	Password    string
	Subdir      *string
	Permissions string
	Crud        *CrudType
}

// Cors contains settings related to Cross-Origin Resource Sharing (CORS)
type Cors struct {
	Origin      string
	Credentials bool
}

// ParseConfig parses the application configuration an sets defaults.
func ParseConfig(path string) *Config {
	// Initialize and log configuration loading
	var cfg = &Config{}
	log.WithField("path", path).Debug("Parsing config file")
	//setDefaults() // Apply default configuration values
	// Determine configuration file location
	if path != "" {
		viper.SetConfigFile(path) // Use provided path
	} else {
		viper.SetConfigName("config")       // Search for default file name
		viper.AddConfigPath("./config")     // Add local config directory
		viper.AddConfigPath("$HOME/.swd")   // Check user's Switcher directory
		viper.AddConfigPath("$HOME/.david") // Check user's David directory
		viper.AddConfigPath(".")            // Include current directory
	}
	// Read and validate configuration file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %s", err)) // Propagate error with details
	}
	err = viper.Unmarshal(&cfg) // Unmarshall values into Config struct
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error parsing config file: %s", err)) // Propagate error with context
	}
	log.WithField("path", viper.ConfigFileUsed()).Debug("Finished Unmarshalling config file")

	// Set production mode for logging in NDJSON format
	cfg.Log.Production = viper.GetBool("Log.Production")
	cfg.Log.Debug = viper.GetBool("Log.Debug")

	// Process user permissions
	for user := range viper.GetStringMap("Users") {
		log.WithField("user", user).Debug("Processing user permissions") // Log user permissions processing
		if cfg.Users[user] == nil {
			log.WithField("user", user).Error("User not found in config file") // Log error with context
			log.WithError(errors.New("cannot launch David without a defined user")).Error("user: " + user + " is not defined in the config file")
			os.Exit(65)
		}
		permissions := viper.GetString(fmt.Sprintf("Users.%s.permissions", user)) // Access specific user permissions
		cfg.Users[user].Crud = &CrudType{Crud: permissions}                       // Set user's CRUD permissions object
		err := FormatCrud(context.Background(), user, cfg)                        // Further process and validate permissions
		if err != nil {
			log.WithError(err).WithField("user", user).Error("Error parsing crud string from config file") // log error with context
		}
		log.WithFields(logrus.Fields{"user": user,
			"crud": cfg.Users[user].Crud}).Debug("Parsed crud string from config file") // Log parsed permissions
	}

	// Validate TLS configuration (if present)
	if cfg.TLS != nil {
		if _, err := os.Stat(cfg.TLS.KeyFile); err != nil {
			log.Fatal(fmt.Errorf("TLS keyFile doesn't exist: %s", err)) // Check for and log missing key file error
		}
		if _, err := os.Stat(cfg.TLS.CertFile); err != nil {
			log.Fatal(fmt.Errorf("TLS certFile doesn't exist: %s", err)) // Check for and log missing cert file error
		}
	}
	// Enable config hot reload and update
	viper.WatchConfig()
	// Register callback for handling config changes
	viper.OnConfigChange(cfg.handleConfigUpdate)
	// Create base and user directories if necessary
	cfg.createBaseAndUserDirectoriesIfNeeded()
	// Return successfully parsed configuration
	return cfg
}

// AuthenticationNeeded returns whether users are defined and authentication is required
func (cfg *Config) AuthenticationNeeded() bool {
	return cfg.Users != nil && len(cfg.Users) != 0
}

func (cfg *Config) handleConfigUpdate(e fsnotify.Event) {
	// Recover from any panics during config update
	defer func() {
		r := recover()
		switch t := r.(type) {
		case string:
			log.WithError(errors.New(t)).Error("Error updating configuration. Please restart the server...")
		case error:
			log.WithError(t).Error("Error updating configuration. Please restart the server...")
		}
	}()

	// Log the config file change
	log.WithField("path", e.Name).Debug("Config file changed")

	// Open the config file for reading
	file, err := os.Open(e.Name)
	if err != nil {
		log.WithField("path", e.Name).Warn("Error reloading config")
	}

	// Create a new Config object to hold updated values
	var updatedCfg = &Config{}

	// Read the config file into the viper instance
	viper.ReadConfig(file)
	// Unmarshal the viper config into the updatedCfg object
	if err := viper.Unmarshal(updatedCfg); err != nil {
		log.WithError(err).Error("Error parsing config file")
		return
	}
	updateConfig(cfg, updatedCfg)
}

// Call the updateConfig function to merge changes
func updateConfig(cfg *Config, updatedCfg *Config) {
	for username := range cfg.Users {
		if updatedCfg.Users[username] == nil {
			log.WithField("user", username).Debug("Removed User from configuration")
			delete(cfg.Users, username)
		}
	}
	// Process added and updated users
	for username, userInformationChange := range updatedCfg.Users {
		if cfg.Users[username] == nil {
			log.WithField("user", username).Info("Added User to configuration")
			cfg.Users[username] = userInformationChange
		} else {
			// Update password, subdir, and crud if changed
			if cfg.Users[username].Password != userInformationChange.Password {
				log.WithField("user", username).Info("Updated password of user")
				cfg.Users[username].Password = userInformationChange.Password
			}
			if cfg.Users[username].Subdir != userInformationChange.Subdir {
				log.WithField("user", username).Info("Updated subdir of user")
				cfg.Users[username].Subdir = userInformationChange.Subdir
			}
			if cfg.Users[username].Crud != userInformationChange.Crud {
				cfg.Users[username].Crud = &CrudType{Crud: userInformationChange.Permissions}
				err := FormatCrud(context.Background(), username, cfg)
				if err != nil {
					log.WithError(err).WithField("user", username).Error("Error parsing crud string from config file")
				}
				log.WithField("user", username).Info("Updated crud of user")
			}
		}
	}
	// Update base and user directories if needed
	cfg.createBaseAndUserDirectoriesIfNeeded()

	// Update logging settings
	// Log.Production should never be updated during actual production, therefore it's not included here
	if cfg.Log.Debug != updatedCfg.Log.Debug {
		cfg.Log.Debug = updatedCfg.Log.Debug
		log.WithField("enabled", cfg.Log.Debug).Debug("Set debug mode")
	}
	if cfg.Log.Create != updatedCfg.Log.Create {
		cfg.Log.Create = updatedCfg.Log.Create
		log.WithField("enabled", cfg.Log.Create).Debug("Set logging for create operations")
	}
	if cfg.Log.Read != updatedCfg.Log.Read {
		cfg.Log.Read = updatedCfg.Log.Read
		log.WithField("enabled", cfg.Log.Read).Debug("Set logging for read operations")
	}
	if cfg.Log.Update != updatedCfg.Log.Update {
		cfg.Log.Update = updatedCfg.Log.Update
		log.WithField("enabled", cfg.Log.Update).Debug("Set logging for update operations")
	}
	if cfg.Log.Delete != updatedCfg.Log.Delete {
		cfg.Log.Delete = updatedCfg.Log.Delete
		log.WithField("enabled", cfg.Log.Delete).Debug("Set logging for delete operations")
	}
}

// createBaseAndUserDirectoriesIfNeeded creates the base directory and individual
// user directories if they don't already exist.
func (cfg *Config) createBaseAndUserDirectoriesIfNeeded() {
	// Check if the base directory already exists.
	if _, err := os.Stat(cfg.Dir); os.IsNotExist(err) {
		mkdirErr := os.Mkdir(cfg.Dir, os.ModePerm)
		if mkdirErr != nil {
			log.WithField("path", cfg.Dir).WithField("error", err).Warn("Can't create base dir")
			return
		}
		log.WithField("path", cfg.Dir).Info("Created base dir")
	}

	// Create individual user directories if they have a defined subdirectory.
	for _, user := range cfg.Users {
		if user.Subdir != nil {
			path := filepath.Join(cfg.Dir, *user.Subdir) // Use path.Join directly for clarity.
			_, pathErr := os.Stat(path)
			if os.IsNotExist(pathErr) {
				os.Mkdir(path, os.ModePerm)
				log.WithField("path", path).Info("Created user dir")
			}
		}
	}
}
