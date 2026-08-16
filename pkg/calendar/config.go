package calendar

import "time"

type CalendarProvider string

const (
	ProviderGoogleCalendar CalendarProvider = "google_calendar"
	ProviderOutlook        CalendarProvider = "outlook_office365"
	ProviderCalCom         CalendarProvider = "cal_com"
	ProviderCalendly       CalendarProvider = "calendly"
	ProviderCustomWebhook  CalendarProvider = "custom_webhook"
)

type CalendarConfig struct {
	Provider CalendarProvider `json:"provider"`

	// 1. Google Calendar / Outlook Credentials
	APIKey      string `json:"api_key,omitempty"`
	OAuthToken  string `json:"oauth_token,omitempty"`
	CalendarID  string `json:"calendar_id,omitempty"` // Primary calendar or specific ID

	// 2. Cal.com / Calendly Specifics
	EventTypeID string `json:"event_type_id,omitempty"`
	UserURI     string `json:"user_uri,omitempty"`

	// 3. Custom Webhook
	WebhookURL string `json:"webhook_url,omitempty"`
	AuthHeader string `json:"auth_header,omitempty"`

	// Timezone
	TimeZone string `json:"timezone,omitempty"` // e.g. "Europe/Paris"
}

type TimeSlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Available bool      `json:"available"`
}

type AvailabilityRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Duration  int       `json:"duration_minutes"` // in minutes, e.g. 30
}

type AvailabilityResponse struct {
	Slots    []TimeSlot `json:"slots"`
	TimeZone string     `json:"timezone"`
}

type BookingRequest struct {
	Slot        TimeSlot          `json:"slot"`
	Attendee    AttendeeInfo      `json:"attendee"`
	Summary     string            `json:"summary,omitempty"`
	Description string            `json:"description,omitempty"`
	CustomVars  map[string]string `json:"custom_vars,omitempty"`
}

type AttendeeInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type BookingResponse struct {
	BookingID string    `json:"booking_id"`
	Status    string    `json:"status"` // "CONFIRMED", "PENDING", "FAILED"
	HtmlLink  string    `json:"html_link,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
