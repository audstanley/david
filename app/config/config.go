// Package config provides configuration management using Viper
package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the david server
type Config struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
	Dir     string `mapstructure:"dir"`
	Prefix  string `mapstructure:"prefix"`

	Database   DatabaseConfig        `mapstructure:"database"`
	WebDAV     WebDAVConfig          `mapstructure:"webdav"`
	API        APIConfig             `mapstructure:"api"`
	Auth       AuthConfig            `mapstructure:"auth"`
	Recurrence RecurrenceConfig      `mapstructure:"recurrence"`
	Log        LogConfig             `mapstructure:"log"`
	Admin      AdminConfig           `mapstructure:"admin"`
	TLS        *TLSConfig            `mapstructure:"tls"`
	Cors       CorsConfig            `mapstructure:"cors"`
	Users      map[string]UserConfig `mapstructure:"users"`
	Hash       HashConfig            `mapstructure:"hash"`
}

// DatabaseConfig holds LevelDB database configuration
type DatabaseConfig struct {
	BaseDir        string `mapstructure:"base_dir"`
	CalendarsDB    string `mapstructure:"calendars_db"`
	EventsDB       string `mapstructure:"events_db"`
	TodosDB        string `mapstructure:"todos_db"`
	JournalsDB     string `mapstructure:"journals_db"`
	FreeBusyDB     string `mapstructure:"freebusy_db"`
	TimezonesDB    string `mapstructure:"timezones_db"`
	UsersDB        string `mapstructure:"users_db"`
	JWTBlacklistDB string `mapstructure:"jwt_blacklist_db"`
	AuditDB        string `mapstructure:"audit_db"`
}

// WebDAVConfig holds WebDAV server configuration
type WebDAVConfig struct {
	Address string     `mapstructure:"address"`
	Port    string     `mapstructure:"port"`
	Prefix  string     `mapstructure:"prefix"`
	TLS     *TLSConfig `mapstructure:"tls"`
}

// APIConfig holds REST API configuration
type APIConfig struct {
	Address      string          `mapstructure:"address"`
	Port         string          `mapstructure:"port"`
	CORSEnabled  bool            `mapstructure:"cors_enabled"`
	CORSOrigins  []string        `mapstructure:"cors_origins"`
	AuthMethods  []string        `mapstructure:"auth_methods"`
	RateLimiting RateLimitConfig `mapstructure:"rate_limiting"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	DefaultLimit int    `mapstructure:"default_limit"`
	Window       string `mapstructure:"window"`
	PerIP        bool   `mapstructure:"per_ip"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWT    JWTConfig    `mapstructure:"jwt"`
	APIKey APIKeyConfig `mapstructure:"api_key"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	Expiry        string `mapstructure:"expiry"`
	RefreshExpiry string `mapstructure:"refresh_expiry"`
}

// APIKeyConfig holds API key configuration
type APIKeyConfig struct {
	Prefix string `mapstructure:"prefix"`
	Expiry string `mapstructure:"expiry"`
}

// RecurrenceConfig holds recurrence cache configuration
type RecurrenceConfig struct {
	CacheExpiryDays   int `mapstructure:"cache_expiry_days"`
	MaxCacheInstances int `mapstructure:"max_cache_instances"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string            `mapstructure:"level"`
	Production bool              `mapstructure:"production"`
	Rotation   LogRotationConfig `mapstructure:"rotation"`
	Debug      bool              `mapstructure:"debug"`
}

// LogRotationConfig holds log rotation settings
type LogRotationConfig struct {
	MaxSize    int `mapstructure:"max_size"`
	MaxBackups int `mapstructure:"max_backups"`
	MaxAge     int `mapstructure:"max_age"`
}

// AdminConfig holds admin feature configuration
type AdminConfig struct {
	Statistics AdminStatisticsConfig `mapstructure:"statistics"`
	Audit      AdminAuditConfig      `mapstructure:"audit"`
}

// AdminStatisticsConfig holds statistics configuration
type AdminStatisticsConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	CacheDuration string `mapstructure:"cache_duration"`
}

// AdminAuditConfig holds audit configuration
type AdminAuditConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	RetentionDays int  `mapstructure:"retention_days"`
}

// TLSConfig holds TLS/SSL configuration
type TLSConfig struct {
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// CorsConfig holds CORS configuration
type CorsConfig struct {
	Origin      string `mapstructure:"origin"`
	Credentials bool   `mapstructure:"credentials"`
}

// UserConfig holds user-specific configuration
type UserConfig struct {
	Password    string     `mapstructure:"password"`
	Subdir      string     `mapstructure:"subdir"`
	Permissions string     `mapstructure:"permissions"`
	Crud        CrudConfig `mapstructure:"crud"`
	Role        string     `mapstructure:"role"`
	Email       string     `mapstructure:"email"`
}

// CrudConfig holds CRUD permissions
type CrudConfig struct {
	Create bool `mapstructure:"create"`
	Read   bool `mapstructure:"read"`
	Update bool `mapstructure:"update"`
	Delete bool `mapstructure:"delete"`
}

// HashConfig holds password hashing configuration
type HashConfig struct {
	Algorithm string                 `mapstructure:"algorithm"`
	Params    map[string]interface{} `mapstructure:"params"`
}

// ParseConfig loads configuration from file and environment variables
func ParseConfig(configPath string) *Config {
	v := viper.New()

	// Default configuration
	v.SetDefault("address", "127.0.0.1")
	v.SetDefault("port", "8000")
	v.SetDefault("dir", "/var/lib/david")
	v.SetDefault("prefix", "/webdav")
	v.SetDefault("database.base_dir", "/var/lib/david")
	v.SetDefault("database.calendars_db", "calendars.db")
	v.SetDefault("database.events_db", "events.db")
	v.SetDefault("database.todos_db", "todos.db")
	v.SetDefault("database.journals_db", "journals.db")
	v.SetDefault("database.freebusy_db", "freebusy.db")
	v.SetDefault("database.timezones_db", "timezones.db")
	v.SetDefault("database.users_db", "users.db")
	v.SetDefault("database.jwt_blacklist_db", "jwt_blacklist.db")
	v.SetDefault("database.audit_db", "audit.db")
	v.SetDefault("webdav.address", "127.0.0.1")
	v.SetDefault("webdav.port", "8000")
	v.SetDefault("webdav.prefix", "/webdav")
	v.SetDefault("api.address", "127.0.0.1")
	v.SetDefault("api.port", "8080")
	v.SetDefault("api.cors_enabled", true)
	v.SetDefault("api.auth_methods", []string{"jwt", "api_key", "basic"})
	v.SetDefault("api.rate_limiting.enabled", true)
	v.SetDefault("api.rate_limiting.default_limit", 100)
	v.SetDefault("api.rate_limiting.window", "1m")
	v.SetDefault("api.rate_limiting.per_ip", true)
	v.SetDefault("auth.jwt.expiry", "24h")
	v.SetDefault("auth.jwt.refresh_expiry", "7d")
	v.SetDefault("auth.api_key.prefix", "david_")
	v.SetDefault("auth.api_key.expiry", "365d")
	v.SetDefault("recurrence.cache_expiry_days", 30)
	v.SetDefault("recurrence.max_cache_instances", 1000)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.production", false)
	v.SetDefault("log.rotation.max_size", 100)
	v.SetDefault("log.rotation.max_backups", 5)
	v.SetDefault("log.rotation.max_age", 30)
	v.SetDefault("admin.statistics.enabled", true)
	v.SetDefault("admin.statistics.cache_duration", "5m")
	v.SetDefault("admin.audit.enabled", true)
	v.SetDefault("admin.audit.retention_days", 90)

	// Set config file path if provided
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Auto-detect config file
		configPaths := []string{
			"./config/config.yaml",
			"$HOME/.swd/config.yaml",
			"$HOME/.david/config.yaml",
			"./config.yaml",
		}
		for _, path := range configPaths {
			path = os.ExpandEnv(path)
			if _, err := os.Stat(path); err == nil {
				v.SetConfigFile(path)
				break
			}
		}
	}

	// Bind environment variables
	v.SetEnvPrefix("DAVID")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		// Config file not found, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, proceed with defaults
		} else {
			// Other error reading config
			panic(err)
		}
	}

	// Unmarshal into Config struct
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		panic(err)
	}

	// Expand environment variables in paths
	cfg.Dir = os.ExpandEnv(cfg.Dir)
	cfg.Database.BaseDir = os.ExpandEnv(cfg.Database.BaseDir)
	if cfg.WebDAV.TLS != nil {
		cfg.WebDAV.TLS.CertFile = os.ExpandEnv(cfg.WebDAV.TLS.CertFile)
		cfg.WebDAV.TLS.KeyFile = os.ExpandEnv(cfg.WebDAV.TLS.KeyFile)
	}

	return cfg
}

// GetDatabasePaths returns the full paths to all database files
func (d *DatabaseConfig) GetDatabasePaths() map[string]string {
	paths := make(map[string]string)
	dbPath := filepath.Join(d.BaseDir, "data")

	paths["calendars"] = filepath.Join(dbPath, d.CalendarsDB)
	paths["events"] = filepath.Join(dbPath, d.EventsDB)
	paths["todos"] = filepath.Join(dbPath, d.TodosDB)
	paths["journals"] = filepath.Join(dbPath, d.JournalsDB)
	paths["freebusy"] = filepath.Join(dbPath, d.FreeBusyDB)
	paths["timezones"] = filepath.Join(dbPath, d.TimezonesDB)
	paths["users"] = filepath.Join(dbPath, d.UsersDB)
	paths["jwt_blacklist"] = filepath.Join(dbPath, d.JWTBlacklistDB)
	paths["audit"] = filepath.Join(dbPath, d.AuditDB)

	return paths
}
