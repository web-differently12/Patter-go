package dto

// Brain DTOs
type BrainChatRequest struct {
	Message    string            `json:"message" binding:"required" example:"Hello, what is my order status?"`
	Channel    string            `json:"channel" example:"whatsapp"` // whatsapp, sms, webchat
	ContextCRM map[string]string `json:"context_crm,omitempty"`
	BrainTier  string            `json:"brain_tier,omitempty" example:"smart"` // "fast", "smart", "expert"
}

type BrainChatResponse struct {
	Reply          string `json:"reply" example:"Your order #1234 is currently in transit."`
	ProcessingTier string `json:"processing_tier" example:"enterprise_smart"`
}

// Voice DTOs
type VoiceOutboundRequest struct {
	TargetNumber string `json:"target_number" binding:"required" example:"+15550199"`
	Prompt       string `json:"prompt" binding:"required" example:"You are a friendly appointment reminder."`
	VoiceTier    string `json:"voice_tier,omitempty" example:"premium_hd"`
}

type VoiceCampaignRequest struct {
	Name         string   `json:"name" binding:"required" example:"Q3 Renewal Campaign"`
	TargetNumbers []string `json:"target_numbers" binding:"required"`
	Prompt       string   `json:"prompt" binding:"required"`
	MinJitterSec int      `json:"min_jitter_sec" example:"5"`  // Anti-spam delay jitter
	MaxJitterSec int      `json:"max_jitter_sec" example:"15"` // Anti-spam delay jitter
}

// WhatsApp DTOs
type WhatsAppSendRequest struct {
	RecipientPhone string `json:"recipient_phone" binding:"required" example:"+15550199"`
	Message        string `json:"message" binding:"required" example:"Hello from WhatsApp!"`
	MediaURL       string `json:"media_url,omitempty"`
}

type WhatsAppCampaignRequest struct {
	CampaignName string   `json:"campaign_name" binding:"required"`
	TemplateName string   `json:"template_name" binding:"required" example:"order_update_hsm"`
	Recipients   []string `json:"recipients" binding:"required"`
}

// Messaging SMS/MMS DTOs
type MessageSendRequest struct {
	ToNumber string `json:"to_number" binding:"required" example:"+15550199"`
	Text     string `json:"text" binding:"required" example:"Your verification code is 4920"`
	MediaURL string `json:"media_url,omitempty"`
}

type MessageCampaignRequest struct {
	Name         string            `json:"name" binding:"required"`
	Recipients   []string          `json:"recipients" binding:"required"`
	TemplateText string            `json:"template_text" binding:"required" example:"Hello {{first_name}}, special offer for {{company}}!"`
	Variables    map[string]string `json:"variables,omitempty"`
}

// Avatar DTOs
type AvatarLiveRequest struct {
	FaceURL string `json:"face_url" binding:"required" example:"https://client.com/avatar.png"`
	VoiceID string `json:"voice_id" example:"custom_voice_1"`
}

type AvatarOfflineRequest struct {
	FaceURL    string `json:"face_url" binding:"required"`
	ScriptText string `json:"script_text" binding:"required"`
	VoiceID    string `json:"voice_id" example:"custom_voice_1"`
}

// Meeting DTOs
type MeetingScheduleRequest struct {
	RoomName string `json:"room_name" binding:"required" example:"strategy-sync"`
	Tier     string `json:"tier" example:"interactive_webrtc"` // "interactive_webrtc", "high_capacity"
}

type MeetingBotRequest struct {
	MeetingURL string `json:"meeting_url" binding:"required" example:"https://meet.client.com/abc-defg-hij"`
	BotName    string `json:"bot_name" example:"AI Executive Assistant"`
}
