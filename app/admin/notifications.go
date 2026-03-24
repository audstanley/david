package admin

import (
	"context"
	"fmt"
	"time"
)

// NotificationType defines the type of notification
type NotificationType string

const (
	// NotificationTypeAlert for system alerts
	NotificationTypeAlert NotificationType = "alert"
	// NotificationTypeWarning for warnings
	NotificationTypeWarning NotificationType = "warning"
	// NotificationTypeInfo for informational messages
	NotificationTypeInfo NotificationType = "info"
	// NotificationTypeSuccess for success messages
	NotificationTypeSuccess NotificationType = "success"
)

// NotificationLevel defines severity level
type NotificationLevel int

const (
	LevelLow NotificationLevel = iota
	LevelMedium
	LevelHigh
	LevelCritical
)

// Notification represents an admin notification
type Notification struct {
	ID         string
	Type       NotificationType
	Level      NotificationLevel
	Title      string
	Message    string
	TargetUser string // Empty for all users
	Created    time.Time
	ReadAt     time.Time
	Metadata   map[string]interface{}
}

// NotificationPlugin is the interface for notification plugins
// Implement this interface to add new notification channels
type NotificationPlugin interface {
	// PluginName returns the name of this plugin
	PluginName() string

	// Send sends a notification using this plugin's channel
	Send(ctx context.Context, notif *Notification) error

	// CanSend checks if this plugin can send notifications
	// (e.g., email plugin checks if email is configured)
	CanSend(notif *Notification) bool

	// Cleanup performs any necessary cleanup
	Cleanup() error
}

// NotificationManager handles notification dispatching
type NotificationManager struct {
	plugins []NotificationPlugin
	storage NotificationStorage
	enabled map[string]bool
}

// NotificationStorage defines how notifications are stored
type NotificationStorage interface {
	Save(ctx context.Context, notif *Notification) error
	GetByID(ctx context.Context, id string) (*Notification, error)
	List(ctx context.Context, opts ListOptions) ([]*Notification, error)
	MarkRead(ctx context.Context, id string, readAt time.Time) error
	Delete(ctx context.Context, id string) error
}

// ListOptions for filtering notifications
type ListOptions struct {
	Types      []NotificationType
	Levels     []NotificationLevel
	TargetUser string
	UnreadOnly bool
	Limit      int
	Offset     int
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(storage NotificationStorage) *NotificationManager {
	return &NotificationManager{
		plugins: make([]NotificationPlugin, 0),
		storage: storage,
		enabled: make(map[string]bool),
	}
}

// RegisterPlugin adds a notification plugin
func (m *NotificationManager) RegisterPlugin(plugin NotificationPlugin) {
	m.plugins = append(m.plugins, plugin)
	m.enabled[plugin.PluginName()] = true
}

// UnregisterPlugin removes a notification plugin
func (m *NotificationManager) UnregisterPlugin(name string) {
	for i, plugin := range m.plugins {
		if plugin.PluginName() == name {
			m.plugins = append(m.plugins[:i], m.plugins[i+1:]...)
			delete(m.enabled, name)
			break
		}
	}
}

// EnablePlugin enables a plugin
func (m *NotificationManager) EnablePlugin(name string) {
	m.enabled[name] = true
}

// DisablePlugin disables a plugin
func (m *NotificationManager) DisablePlugin(name string) {
	m.enabled[name] = false
}

// Send creates and dispatches a notification
func (m *NotificationManager) Send(ctx context.Context, notif *Notification) error {
	if notif.ID == "" {
		notif.ID = fmt.Sprintf("notif-%d", time.Now().UnixNano())
	}
	if notif.Created.IsZero() {
		notif.Created = time.Now()
	}

	// Save to storage
	if err := m.storage.Save(ctx, notif); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	// Dispatch to all enabled plugins that can send
	for _, plugin := range m.plugins {
		if !m.enabled[plugin.PluginName()] {
			continue
		}
		if !plugin.CanSend(notif) {
			continue
		}

		if err := plugin.Send(ctx, notif); err != nil {
			// Log error but continue with other plugins
			// TODO: Add logging
			continue
		}
	}

	return nil
}

// Get retrieves a notification by ID
func (m *NotificationManager) Get(ctx context.Context, id string) (*Notification, error) {
	return m.storage.GetByID(ctx, id)
}

// List returns filtered notifications
func (m *NotificationManager) List(ctx context.Context, opts ListOptions) ([]*Notification, error) {
	return m.storage.List(ctx, opts)
}

// MarkRead marks a notification as read
func (m *NotificationManager) MarkRead(ctx context.Context, id string) error {
	return m.storage.MarkRead(ctx, id, time.Now())
}

// Alert creates an alert notification
func (m *NotificationManager) Alert(ctx context.Context, title, message string) error {
	return m.Send(ctx, &Notification{
		Type:    NotificationTypeAlert,
		Level:   LevelHigh,
		Title:   title,
		Message: message,
	})
}

// Warning creates a warning notification
func (m *NotificationManager) Warning(ctx context.Context, title, message string) error {
	return m.Send(ctx, &Notification{
		Type:    NotificationTypeWarning,
		Level:   LevelMedium,
		Title:   title,
		Message: message,
	})
}

// Info creates an info notification
func (m *NotificationManager) Info(ctx context.Context, title, message string) error {
	return m.Send(ctx, &Notification{
		Type:    NotificationTypeInfo,
		Level:   LevelLow,
		Title:   title,
		Message: message,
	})
}

// SystemAlert creates a system-level alert
func (m *NotificationManager) SystemAlert(ctx context.Context, alertType string, data map[string]interface{}) error {
	notif := &Notification{
		Type:     NotificationTypeAlert,
		Level:    LevelHigh,
		Title:    fmt.Sprintf("System %s", alertType),
		Message:  fmt.Sprintf("System alert: %s", alertType),
		Metadata: data,
	}
	return m.Send(ctx, notif)
}

// CapacityWarning creates a storage capacity warning
func (m *NotificationManager) CapacityWarning(ctx context.Context, percentage float64, limit int64, used int64) error {
	notif := &Notification{
		Type:    NotificationTypeWarning,
		Level:   LevelMedium,
		Title:   "Storage Capacity Warning",
		Message: fmt.Sprintf("Storage usage at %.1f%% (%d/%d bytes)", percentage, used, limit),
		Metadata: map[string]interface{}{
			"percentage": percentage,
			"limit":      limit,
			"used":       used,
		},
	}
	return m.Send(ctx, notif)
}

// LoginFailureAlert creates an alert for failed login attempts
func (m *NotificationManager) LoginFailureAlert(ctx context.Context, username, ip string, attempts int) error {
	notif := &Notification{
		Type:    NotificationTypeAlert,
		Level:   LevelHigh,
		Title:   "Multiple Login Failures",
		Message: fmt.Sprintf("User %s failed to login %d times from %s", username, attempts, ip),
		Metadata: map[string]interface{}{
			"username": username,
			"ip":       ip,
			"attempts": attempts,
		},
	}
	return m.Send(ctx, notif)
}

// AuditAnomaly creates an alert for audit anomalies
func (m *NotificationManager) AuditAnomaly(ctx context.Context, anomalyType string, details map[string]interface{}) error {
	notif := &Notification{
		Type:     NotificationTypeAlert,
		Level:    LevelCritical,
		Title:    "Audit Anomaly Detected",
		Message:  fmt.Sprintf("Anomaly detected: %s", anomalyType),
		Metadata: details,
	}
	return m.Send(ctx, notif)
}

// Cleanup cleans up all plugins
func (m *NotificationManager) Cleanup() error {
	for _, plugin := range m.plugins {
		if err := plugin.Cleanup(); err != nil {
			return err
		}
	}
	return nil
}
