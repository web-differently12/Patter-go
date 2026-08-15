package dto

// VoiceOptionsDTO contains agnostic, vendor-neutral voice parameters
type VoiceOptionsDTO struct {
	SilenceTimeoutMs        int     `json:"silence_timeout_ms" example:"5000"`
	MaxDurationSeconds      int     `json:"max_duration_seconds" example:"1800"`
	BackgroundDenoising     bool    `json:"background_denoising" example:"true"`
	BackchannelingEnabled   bool    `json:"backchanneling_enabled" example:"true"`
	InterruptionThresholdMs int     `json:"interruption_threshold_ms" example:"200"`
	VoicemailDetection      bool    `json:"voicemail_detection" example:"true"`
	VADSensitivity          float64 `json:"vad_sensitivity" example:"0.25"`
	EmotionExpressions      bool    `json:"emotion_expressions" example:"true"`
	VoiceQualityTier        string  `json:"voice_quality_tier" example:"premium_hd"` // "standard", "premium_hd", "ultra_natural"
}

// ClientMeetingOptionsDTO contains customizable options for clients to configure meetings without vendor leaks
type ClientMeetingOptionsDTO struct {
	CustomLogoURL         string `json:"custom_logo_url,omitempty" example:"https://client.com/logo.png"`
	WatermarkDisabled     bool   `json:"watermark_disabled" example:"true"`
	LayoutMode            string `json:"layout_mode" example:"grid"` // grid, speaker, spotlight
	RecordingEnabled      bool   `json:"recording_enabled" example:"true"`
	WaitingRoomEnabled     bool   `json:"waiting_room_enabled" example:"false"`
	MaxParticipants       int    `json:"max_participants" example:"50"`
	TranscriptionLanguage string `json:"transcription_language" example:"en-US"`
	CustomDomain          string `json:"custom_domain,omitempty" example:"meet.client.com"`
}

type VoiceCallRequest struct {
	TargetNumber string          `json:"target_number" binding:"required" example:"+15550199"`
	Prompt       string          `json:"prompt" binding:"required" example:"Appointment confirmation assistant."`
	VoiceOptions VoiceOptionsDTO `json:"voice_options,omitempty"`
}

type ClientMeetingConfigRequest struct {
	RoomName       string                  `json:"room_name" binding:"required" example:"strategy-call"`
	ServiceTier    string                  `json:"service_tier" example:"enterprise_interactive"` // "standard_webrtc", "enterprise_interactive"
	MeetingOptions ClientMeetingOptionsDTO `json:"meeting_options,omitempty"`
}
