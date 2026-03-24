// Package icalendar provides parsing and generation of iCalendar files per RFC 2445
package icalendar

import "time"

// Attendee represents a calendar attendee
type Attendee struct {
	Email         string
	Name          string
	RSVP          bool
	Role          string // REQ-PARTICIPANT, OPT-PARTICIPANT, NON-PARTICIPANT
	PartStat      string // NEEDS-ACTION, ACCEPTED, DECLINED, TENTATIVE, DELEGATED
	Cutype        string // INDIVIDUAL, GROUP, RESOURCE, ROOM
	Dir           string
	SentBy        string
	DelegatedFrom []string
	DelegatedTo   []string
}

// Alarm represents a VALARM component
type Alarm struct {
	Action      string // DISPLAY, EMAIL, PROCEDURE
	Trigger     string // Duration or period
	Duration    time.Duration
	Description string
	Attach      string
	Attendees   []string
	Summary     string
}

// RRule represents a recurrence rule (RRULE)
type RRule struct {
	Freq         string // SECONDLY, MINUTELY, HOURLY, DAILY, WEEKLY, MONTHLY, YEARLY
	Interval     int
	BySecond     []int
	ByMinute     []int
	ByHour       []int
	ByDay        []string // MO, TU, WE, TH, FR, SA, SU, +nMO, etc.
	ByMonthDay   []int
	ByYearDay    []int
	ByWeekNo     []int
	ByMonth      []int
	BySetPos     []int
	Count        int
	Until        time.Time
	ExcludeDates []time.Time // EXDATE
}

// Event represents a VEVENT component
type Event struct {
	UID          string
	CalendarUID  string
	DTStamp      time.Time
	DTStart      time.Time
	DTEnd        time.Time
	Duration     *time.Duration
	Summary      string
	Description  string
	Location     string
	Organizer    string
	Attendees    []Attendee
	Status       string // CONFIRMED, TENTATIVE, CANCELLED
	Class        string // PUBLIC, PRIVATE, CONFIDENTIAL
	Priority     int
	Sequence     int
	Categories   []string
	Resources    []string
	RRule        *RRule
	ExDates      []time.Time
	RDates       []time.Time
	RelatedTo    string
	Geo          *Geo
	URL          string
	Created      time.Time
	LastModified time.Time
	Alarms       []Alarm
}

// Todo represents a VTODO component
type Todo struct {
	UID          string
	CalendarUID  string
	DTStamp      time.Time
	DTStart      time.Time
	DUE          time.Time
	Duration     *time.Duration
	Summary      string
	Description  string
	Location     string
	Organizer    string
	Attendees    []Attendee
	Status       string // NEEDS-ACTION, IN-PROCESS, COMPLETED, CANCELLED
	Priority     int
	Sequence     int
	Percent      int
	Categories   []string
	Resources    []string
	RelatedTo    string
	Geo          *Geo
	URL          string
	Created      time.Time
	LastModified time.Time
	Completed    time.Time
	Alarms       []Alarm
}

// Journal represents a VJOURNAL component
type Journal struct {
	UID          string
	CalendarUID  string
	DTStamp      time.Time
	DTStart      time.Time
	Summary      string
	Description  string
	Organizer    string
	Attendees    []Attendee
	Status       string // DRAFT, FINAL, CANCELLED
	Categories   []string
	RelatedTo    string
	URL          string
	Created      time.Time
	LastModified time.Time
}

// FreeBusy represents a VFREEBUSY component
type FreeBusy struct {
	UID         string
	CalendarUID string
	DTStamp     time.Time
	DTStart     time.Time
	DTEnd       time.Time
	FreeBusy    []FreeBusyBlock
	Organizer   string
	Attendees   []Attendee
	URL         string
}

// FreeBusyBlock represents a free/busy time block
type FreeBusyBlock struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
	FreeBusy string // BUSY, BUSY-UNAVAILABLE, BUSY-TENTATIVE, FREE
}

// TimeZone represents a VTIMEZONE component
type TimeZone struct {
	TZID     string
	LastMod  time.Time
	Standard []TimeZoneStandard
	Daylight []TimeZoneDaylight
}

// TimeZoneStandard represents a STANDARD subcomponent
type TimeZoneStandard struct {
	TZName       string
	DTStart      time.Time
	TZOffsetFrom time.Duration
	TZOffsetTo   time.Duration
	RRule        *RRule
}

// TimeZoneDaylight represents a DAYLIGHT subcomponent
type TimeZoneDaylight struct {
	TZName       string
	DTStart      time.Time
	TZOffsetFrom time.Duration
	TZOffsetTo   time.Duration
	RRule        *RRule
}

// Geo represents geographic coordinates
type Geo struct {
	Lat float64
	Lon float64
}

// AlarmAction represents the action type
type AlarmAction string

const (
	AlarmDisplay   AlarmAction = "DISPLAY"
	AlarmEmail     AlarmAction = "EMAIL"
	AlarmProcedure AlarmAction = "PROCEDURE"
)

// ParticipationStatus represents participation status
type ParticipationStatus string

const (
	PartNeedsAction ParticipationStatus = "NEEDS-ACTION"
	PartAccepted    ParticipationStatus = "ACCEPTED"
	PartDeclined    ParticipationStatus = "DECLINED"
	PartTentative   ParticipationStatus = "TENTATIVE"
	PartDelegated   ParticipationStatus = "DELEGATED"
	PartInProcess   ParticipationStatus = "IN-PROCESS"
	PartCompleted   ParticipationStatus = "COMPLETED"
	PartRejected    ParticipationStatus = "REJECTED"
)

// EventStatus represents event status
type EventStatus string

const (
	EventConfirmed EventStatus = "CONFIRMED"
	EventTentative EventStatus = "TENTATIVE"
	EventCancelled EventStatus = "CANCELLED"
)

// TodoStatus represents todo status
type TodoStatus string

const (
	TodoNeedsAction TodoStatus = "NEEDS-ACTION"
	TodoInProcess   TodoStatus = "IN-PROCESS"
	TodoCompleted   TodoStatus = "COMPLETED"
	TodoCancelled   TodoStatus = "CANCELLED"
)

// JournalStatus represents journal status
type JournalStatus string

const (
	JournalDraft     JournalStatus = "DRAFT"
	JournalFinal     JournalStatus = "FINAL"
	JournalCancelled JournalStatus = "CANCELLED"
)

// Class represents classification
type Class string

const (
	ClassPublic       Class = "PUBLIC"
	ClassPrivate      Class = "PRIVATE"
	ClassConfidential Class = "CONFIDENTIAL"
)

// Transparency represents time transparency
type Transparency string

const (
	TranspOpaque      Transparency = "OPAQUE"
	TranspTransparent Transparency = "TRANSPARENT"
)

// ParticipationRole represents participation role
type ParticipationRole string

const (
	RoleReqParticipant ParticipationRole = "REQ-PARTICIPANT"
	RoleOptParticipant ParticipationRole = "OPT-PARTICIPANT"
	RoleNonParticipant ParticipationRole = "NON-PARTICIPANT"
	RoleChair          ParticipationRole = "CHAIR"
)

// CalendarUserType represents calendar user type
type CalendarUserType string

const (
	CutypeIndividual CalendarUserType = "INDIVIDUAL"
	CutypeGroup      CalendarUserType = "GROUP"
	CutypeResource   CalendarUserType = "RESOURCE"
	CutypeRoom       CalendarUserType = "ROOM"
)

// FreeBusyType represents free/busy type
type FreeBusyType string

const (
	FbTypeBusy            FreeBusyType = "BUSY"
	FbTypeBusyUnavailable FreeBusyType = "BUSY-UNAVAILABLE"
	FbTypeBusyTentative   FreeBusyType = "BUSY-TENTATIVE"
	FbTypeFree            FreeBusyType = "FREE"
)

// Freq represents recurrence frequency
type Freq string

const (
	FreqSecondly Freq = "SECONDLY"
	FreqMinutely Freq = "MINUTELY"
	FreqHourly   Freq = "HOURLY"
	FreqDaily    Freq = "DAILY"
	FreqWeekly   Freq = "WEEKLY"
	FreqMonthly  Freq = "MONTHLY"
	FreqYearly   Freq = "YEARLY"
)
