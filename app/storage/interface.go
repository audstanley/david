package storage

import (
	"time"

	"github.com/audstanley/david/app/icalendar"
)

// User represents a user in the system
type User struct {
	ID           string
	Username     string
	PasswordHash string
	Email        string
	DisplayName  string
	Role         string
	APIKeys      []string
	Created      time.Time
	Updated      time.Time
}

// Calendar represents a calendar
type Calendar struct {
	UID         string
	OwnerID     string
	DisplayName string
	Description string
	Color       string
	IsPublic    bool
	PublicHash  string
	Timezone    string
	Created     time.Time
	Updated     time.Time
}

// Event represents a calendar event
type Event struct {
	UID          string
	CalendarUID  string
	DTStamp      time.Time
	DTStart      time.Time
	DTEnd        time.Time
	Summary      string
	Description  string
	Location     string
	Organizer    string
	Attendees    []icalendar.Attendee
	Status       string
	Class        string
	Priority     int
	Sequence     int
	Categories   []string
	RRule        *icalendar.RRule
	ExDates      []time.Time
	Created      time.Time
	LastModified time.Time
}

// Todo represents a calendar todo
type Todo struct {
	UID          string
	CalendarUID  string
	DTStamp      time.Time
	Due          time.Time
	Summary      string
	Description  string
	Status       string
	Priority     int
	Percent      int
	Sequence     int
	Created      time.Time
	LastModified time.Time
	Completed    time.Time
}

// Journal represents a calendar journal
type Journal struct {
	UID          string
	CalendarUID  string
	DTStamp      time.Time
	DTCreated    time.Time
	Summary      string
	Description  string
	Status       string
	Sequence     int
	Created      time.Time
	LastModified time.Time
}

// FreeBusyBlock represents a free/busy time block
type FreeBusyBlock struct {
	UID    string
	Start  time.Time
	End    time.Time
	Status string // FREE, BUSY, BUSY-UNAVAILABLE
}

// VTimeZone represents a VTIMEZONE component
type VTimeZone struct {
	TZID       string
	Components []icalendar.Component
}

// RecurrenceInstance represents a single instance of a recurring event
type RecurrenceInstance struct {
	EventUID    string
	InstanceUID string
	DTStart     time.Time
	DTEnd       time.Time
}

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID         string
	Timestamp  time.Time
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	Details    string
}

// UserStore defines the interface for user storage operations
type UserStore interface {
	Create(userID, username, passwordHash, email, displayName string, role string) error
	GetByID(userID string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(userID string, updates map[string]interface{}) error
	Delete(userID string) error
	List(limit, offset int) ([]*User, error)
}

// CalendarStore defines the interface for calendar storage operations
type CalendarStore interface {
	Create(calendarUID, ownerID, displayName, description, color, timezone string, isPublic bool) error
	GetByUID(calendarUID string) (*Calendar, error)
	GetByName(ownerID, name string) (*Calendar, error)
	Update(calendarUID string, updates map[string]interface{}) error
	Delete(calendarUID string) error
	List(ownerID string, limit, offset int) ([]*Calendar, error)
	IsPublic(calendarUID string) (bool, error)
}

// EventStore defines the interface for event storage operations
type EventStore interface {
	Create(event *Event) error
	GetByUID(eventUID string) (*Event, error)
	Update(eventUID string, event *Event, sequence int) error
	Delete(eventUID string) error
	List(calendarUID string, start, end time.Time, limit, offset int) ([]*Event, error)
	GetInDateRange(calendarUID, eventUID string, start, end time.Time) (*Event, error)
}

// TodoStore defines the interface for todo storage operations
type TodoStore interface {
	Create(todo *Todo) error
	GetByUID(todoUID string) (*Todo, error)
	Update(todoUID string, todo *Todo, sequence int) error
	Delete(todoUID string) error
	List(calendarUID string, limit, offset int) ([]*Todo, error)
}

// JournalStore defines the interface for journal storage operations
type JournalStore interface {
	Create(journal *Journal) error
	GetByUID(journalUID string) (*Journal, error)
	Update(journalUID string, journal *Journal) error
	Delete(journalUID string) error
	List(calendarUID string, limit, offset int) ([]*Journal, error)
}

// FreeBusyStore defines the interface for free/busy storage operations
type FreeBusyStore interface {
	Generate(calendarUIDs []string, start, end time.Time) ([]FreeBusyBlock, error)
	Get(calendarUID string, start, end time.Time) ([]FreeBusyBlock, error)
}

// TimeZoneStore defines the interface for timezone storage operations
type TimeZoneStore interface {
	Store(tzID string, vtimezone *VTimeZone) error
	Get(tzID string) (*VTimeZone, error)
	Update(tzID string, vtimezone *VTimeZone) error
	Delete(tzID string) error
}

// RecurrenceStore defines the interface for recurrence storage operations
type RecurrenceStore interface {
	Store(eventUID string, instances []RecurrenceInstance) error
	Get(eventUID string, start, end time.Time) ([]RecurrenceInstance, error)
	Delete(eventUID string) error
}

// AuditStore defines the interface for audit log storage operations
type AuditStore interface {
	Log(userID, action, entityType, entityID, details string) error
	GetByUser(userID string, limit, offset int) ([]*AuditEntry, error)
	GetByEntity(entityType, entityID string, limit, offset int) ([]*AuditEntry, error)
	GetByDateRange(start, end time.Time, limit, offset int) ([]*AuditEntry, error)
}

// ShareManager defines the interface for calendar share operations
type ShareManager interface {
	Create(calendarUID, targetUserID, role string) error
	GetByCalendar(calendarUID string) ([]*CalendarShare, error)
	Delete(calendarUID, targetUserID string) error
	CheckAccess(calendarUID, userID string) (string, error)
	ListAll(calendarUID string) ([]*CalendarShare, error)
}

// CalendarShare defines a calendar share entry
type CalendarShare struct {
	ID           string
	CalendarUID  string
	TargetUserID string
	Role         string
	Created      time.Time
	CreatedBy    string
}

// Share roles
const (
	ShareRoleRead     = "read"
	ShareRoleWrite    = "write"
	ShareRoleAdmin    = "admin"
	ShareRoleEveryone = "everyone"
)
