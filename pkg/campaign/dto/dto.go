package dto

import "time"

type MessageType string

const (
	TypeFreeText MessageType = "TEXT"
	TypeImage    MessageType = "IMAGE"
	TypeVideo    MessageType = "VIDEO"
	TypeDocument MessageType = "DOCUMENT"
	TypeVoice    MessageType = "VOICE_CALL"
)

type AudienceFilter struct {
	MatchLogic         string   `json:"match_logic"`          // "AND" | "OR"
	ContactCategoryIDs []string `json:"contact_category_ids"` // Categories to include
	RequiredTags       []string `json:"required_tags"`        // Tags required
	ExcludedTags       []string `json:"excluded_tags"`        // Tags excluded
}

type TargetContact struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Phone       string            `json:"phone"` // E.164 format (e.g., +15551234567)
	CategoryIDs []string          `json:"category_ids,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	OptedOut    bool              `json:"opted_out"`
	Variables   map[string]string `json:"variables,omitempty"`
}

type CampaignConfig struct {
	CampaignID       string         `json:"campaign_id,omitempty"`
	TenantID         string         `json:"tenant_id,omitempty"`
	Name             string         `json:"name" binding:"required"`
	Channel          string         `json:"channel" binding:"required"` // "WHATSAPP", "SMS", "VOICE"
	SessionNames     []string       `json:"session_names" binding:"required"`
	MessageType      MessageType    `json:"message_type"`
	MessageContent   string         `json:"message_content"`
	MediaURLs        []string       `json:"media_urls,omitempty"`
	RandomDelaySec   int            `json:"random_delay_sec"` // Anti-spam jitter base delay in seconds
	StartImmediately bool           `json:"start_immediately"`
	ScheduledFor     *time.Time     `json:"scheduled_for,omitempty"`
	AudienceFilter   AudienceFilter `json:"audience_filter"`
	Targets          []TargetContact `json:"targets,omitempty"`
}

type VoiceCampaignRequest struct {
	CampaignName     string   `json:"campaign_name" binding:"required"`
	FromNumbers      []string `json:"from_numbers" binding:"required"`
	PromptTemplate   string   `json:"prompt_template" binding:"required"`
	RandomDelaySec   int      `json:"random_delay_sec"`
	TargetPhoneNumbers []string `json:"target_phone_numbers" binding:"required"`
}

type CampaignResponse struct {
	CampaignID     string            `json:"campaign_id"`
	TenantID       string            `json:"tenant_id"`
	Name           string            `json:"name"`
	Channel        string            `json:"channel"`
	Status         string            `json:"status"` // "DRAFT", "RUNNING", "PAUSED", "COMPLETED", "FAILED"
	TotalContacts  int               `json:"total_contacts"`
	SentCount      int               `json:"sent_count"`
	DeliveredCount int               `json:"delivered_count"`
	ReadCount      int               `json:"read_count"`
	FailedCount    int               `json:"failed_count"`
	CreatedAt      time.Time         `json:"created_at"`
	ScheduledFor   *time.Time        `json:"scheduled_for,omitempty"`
}

type CampaignStatusResponse struct {
	CampaignID     string    `json:"campaign_id"`
	Status         string    `json:"status"`
	TotalContacts  int       `json:"total_contacts"`
	SentCount      int       `json:"sent_count"`
	DeliveredCount int       `json:"delivered_count"`
	ReadCount      int       `json:"read_count"`
	FailedCount    int       `json:"failed_count"`
	Progress       float64   `json:"progress_percent"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
