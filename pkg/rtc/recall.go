package rtc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/core"
)

type MeetingPlatform string

const (
	PlatformZoom            MeetingPlatform = "zoom"
	PlatformGoogleMeet      MeetingPlatform = "google_meet"
	PlatformMicrosoftTeams MeetingPlatform = "teams"
	PlatformWebex          MeetingPlatform = "webex"
)

type RecordingMode string

const (
	RecordingSpeakerView RecordingMode = "speaker_view"
	RecordingGalleryView RecordingMode = "gallery_view"
	RecordingAudioOnly   RecordingMode = "audio_only"
)

type BotStatus string

const (
	StatusJoining   BotStatus = "joining"
	StatusInCall    BotStatus = "in_call"
	StatusRecording BotStatus = "recording"
	StatusLeft      BotStatus = "left"
	StatusFailed    BotStatus = "failed"
)

type TranscriptionOptions struct {
	Provider         string   `json:"provider,omitempty"`       // "assemblyai", "deepgram", "gladia", "whisper"
	Language         string   `json:"language,omitempty"`       // "fr", "en", "es", "de", "auto"
	CustomVocabulary []string `json:"custom_vocabulary,omitempty"`
}

type AutomaticLeaveOptions struct {
	EveryoneLeftTimeoutSec int  `json:"everyone_left_timeout_sec,omitempty"` // Default 60s
	SilenceTimeoutSec      int  `json:"silence_timeout_sec,omitempty"`       // Default 300s
	MaxDurationMinutes     int  `json:"max_duration_minutes,omitempty"`      // Default 180m
	LeaveWhenHostLeaves    bool `json:"leave_when_host_leaves"`
}

type VideoOptions struct {
	Layout            string `json:"layout,omitempty"`              // "speaker", "gallery"
	Resolution        string `json:"resolution,omitempty"`          // "720p", "1080p"
	WatermarkImageURL string `json:"watermark_image_url,omitempty"` // White-label logo
}

type TenantMeetingSettings struct {
	TenantID             string                `json:"tenant_id"`
	DefaultBotName       string                `json:"default_bot_name"`
	DefaultAvatarURL     string                `json:"default_avatar_url"`
	DefaultLanguage      string                `json:"default_language"`
	TranscriptionOptions TranscriptionOptions  `json:"transcription_options"`
	AutomaticLeave       AutomaticLeaveOptions `json:"automatic_leave"`
	VideoOptions         VideoOptions          `json:"video_options"`
	WebhookEvents        []string              `json:"webhook_events"` // "bot.status_change", "transcript.chunk", "recording.done"
	WebhookURL           string                `json:"webhook_url"`
}

type MeetingProfile struct {
	ProfileID                     string `json:"profile_id"`
	TenantID                      string `json:"tenant_id"`
	Name                          string `json:"name" binding:"required"`
	EnableSpeakerDiarization      bool   `json:"enable_speaker_diarization"`
	EnableActionItemsExtraction bool   `json:"enable_action_items_extraction"`
	EnableParticipantSentiment    bool   `json:"enable_participant_sentiment"`
	EnableLiveTranslation         bool   `json:"enable_live_translation"`
	TargetTranslationLanguage     string `json:"target_translation_language,omitempty"`
	EnableScreenShareRecording    bool   `json:"enable_screen_share_recording"`
	SummaryTemplate               string `json:"summary_template"`
	AutoLeaveOnSilenceMinutes     int    `json:"auto_leave_on_silence_minutes"`
	AutoLeaveWhenEveryoneLeft     bool   `json:"auto_leave_when_everyone_left"`
	CustomAvatarVideoURL          string `json:"custom_avatar_video_url,omitempty"`
}

type CreateMeetingBotRequest struct {
	MeetingURL                string                 `json:"meeting_url" binding:"required"`
	ProfileID                 string                 `json:"profile_id,omitempty"`
	BotName                   string                 `json:"bot_name,omitempty"`
	AvatarURL                 string                 `json:"avatar_url,omitempty"`
	Platform                  MeetingPlatform        `json:"platform,omitempty"`
	RecordingMode             RecordingMode          `json:"recording_mode,omitempty"`
	EnableRealtimeTranscript  bool                   `json:"enable_realtime_transcript"`
	EnableRealtimeAudioStream bool                   `json:"enable_realtime_audio_stream"`
	EnableChatMessaging       bool                   `json:"enable_chat_messaging"`
	TranscriptionOptions      *TranscriptionOptions  `json:"transcription_options,omitempty"`
	AutomaticLeave            *AutomaticLeaveOptions `json:"automatic_leave,omitempty"`
	VideoOptions              *VideoOptions          `json:"video_options,omitempty"`
	Language                  string                 `json:"language,omitempty"`
	WebhookURL                string                 `json:"webhook_url,omitempty"`
	Metadata                  map[string]string      `json:"metadata,omitempty"`
}

type SendMeetingChatMessageRequest struct {
	BotID   string `json:"bot_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type MeetingBotResponse struct {
	BotID                   string                `json:"bot_id"`
	TenantID                string                `json:"tenant_id"`
	ProfileID               string                `json:"profile_id,omitempty"`
	MeetingURL              string                `json:"meeting_url"`
	BotName                 string                `json:"bot_name"`
	AvatarURL               string                `json:"avatar_url,omitempty"`
	Platform                MeetingPlatform       `json:"platform"`
	Status                  BotStatus             `json:"status"`
	Language                string                `json:"language"`
	RecordingMode           RecordingMode         `json:"recording_mode"`
	VideoRecordingURL       string                `json:"video_recording_url,omitempty"`
	TranscriptURL           string                `json:"transcript_url,omitempty"`
	AudioStreamWebSocketURL string                `json:"audio_stream_websocket_url,omitempty"`
	TranscriptionOptions    TranscriptionOptions  `json:"transcription_options"`
	AutomaticLeave          AutomaticLeaveOptions `json:"automatic_leave"`
	CreatedAt               time.Time             `json:"created_at"`
	ProviderUsed            string                `json:"provider_used"` // e.g. "patter_native_rtc", "meetingbaas_engine" or "patter_premium_engine"
}

type MeetingEngineService interface {
	SaveTenantMeetingSettings(ctx context.Context, tenantID string, settings TenantMeetingSettings) (*TenantMeetingSettings, error)
	GetTenantMeetingSettings(ctx context.Context, tenantID string) (*TenantMeetingSettings, error)
	CreateMeetingProfile(ctx context.Context, tenantID string, prof MeetingProfile) (*MeetingProfile, error)
	GetMeetingProfile(ctx context.Context, tenantID, profileID string) (*MeetingProfile, error)
	ListMeetingProfiles(ctx context.Context, tenantID string) ([]*MeetingProfile, error)
	CreateMeetingBot(ctx context.Context, tenantID string, req CreateMeetingBotRequest) (*MeetingBotResponse, error)
	GetBotStatus(ctx context.Context, tenantID, botID string) (*MeetingBotResponse, error)
	ListBots(ctx context.Context, tenantID string) ([]*MeetingBotResponse, error)
	SendChatMessage(ctx context.Context, tenantID string, req SendMeetingChatMessageRequest) error
	LeaveMeeting(ctx context.Context, tenantID, botID string) error
}

type meetingEngineService struct {
	apiKey         string
	costOptimizer  *core.CostOptimizationEngine
	httpClient     *http.Client
	logger         *slog.Logger
	mu             sync.RWMutex
	tenantSettings map[string]*TenantMeetingSettings
	profiles       map[string]*MeetingProfile
	bots           map[string]*MeetingBotResponse
}

func NewMeetingEngineService(apiKey string, logger *slog.Logger) MeetingEngineService {
	if logger == nil {
		logger = core.GetLogger()
	}
	svc := &meetingEngineService{
		apiKey:         apiKey,
		costOptimizer:  core.NewCostOptimizationEngine(logger),
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		logger:         logger.With("component", "patter_meeting_engine"),
		tenantSettings: make(map[string]*TenantMeetingSettings),
		profiles:       make(map[string]*MeetingProfile),
		bots:           make(map[string]*MeetingBotResponse),
	}

	svc.tenantSettings["default_tenant"] = &TenantMeetingSettings{
		TenantID:         "default_tenant",
		DefaultBotName:   "Patter AI Assistant",
		DefaultAvatarURL: "https://patter.ai/assets/avatar_white_label.png",
		DefaultLanguage:  "fr",
		TranscriptionOptions: TranscriptionOptions{
			Provider: "assemblyai",
			Language: "fr",
		},
		AutomaticLeave: AutomaticLeaveOptions{
			EveryoneLeftTimeoutSec: 60,
			SilenceTimeoutSec:      300,
			MaxDurationMinutes:     180,
			LeaveWhenHostLeaves:    true,
		},
		VideoOptions: VideoOptions{
			Layout:     "speaker",
			Resolution: "1080p",
		},
		WebhookEvents: []string{"bot.status_change", "transcript.chunk", "recording.done"},
	}

	return svc
}

func (s *meetingEngineService) SaveTenantMeetingSettings(ctx context.Context, tenantID string, settings TenantMeetingSettings) (*TenantMeetingSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	settings.TenantID = tenantID
	s.tenantSettings[tenantID] = &settings
	return &settings, nil
}

func (s *meetingEngineService) GetTenantMeetingSettings(ctx context.Context, tenantID string) (*TenantMeetingSettings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	settings, ok := s.tenantSettings[tenantID]
	if !ok {
		return &TenantMeetingSettings{
			TenantID:       tenantID,
			DefaultBotName: "Patter AI Agent",
			DefaultLanguage: "fr",
			TranscriptionOptions: TranscriptionOptions{
				Provider: "assemblyai",
				Language: "fr",
			},
			AutomaticLeave: AutomaticLeaveOptions{
				EveryoneLeftTimeoutSec: 60,
				SilenceTimeoutSec:      300,
			},
		}, nil
	}
	return settings, nil
}

func (s *meetingEngineService) CreateMeetingProfile(ctx context.Context, tenantID string, prof MeetingProfile) (*MeetingProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	profID := "mp_" + uuid.New().String()[:8]
	prof.ProfileID = profID
	prof.TenantID = tenantID
	if prof.SummaryTemplate == "" {
		prof.SummaryTemplate = "EXECUTIVE_SUMMARY"
	}

	s.profiles[tenantID+":"+profID] = &prof
	return &prof, nil
}

func (s *meetingEngineService) GetMeetingProfile(ctx context.Context, tenantID, profileID string) (*MeetingProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prof, ok := s.profiles[tenantID+":"+profileID]
	if !ok {
		return nil, fmt.Errorf("meeting profile %s not found", profileID)
	}
	return prof, nil
}

func (s *meetingEngineService) ListMeetingProfiles(ctx context.Context, tenantID string) ([]*MeetingProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*MeetingProfile
	for _, prof := range s.profiles {
		if prof.TenantID == tenantID {
			result = append(result, prof)
		}
	}
	return result, nil
}

func (s *meetingEngineService) CreateMeetingBot(ctx context.Context, tenantID string, req CreateMeetingBotRequest) (*MeetingBotResponse, error) {
	botID := "bot_" + uuid.New().String()[:8]

	tenantCfg, _ := s.GetTenantMeetingSettings(ctx, tenantID)

	botName := req.BotName
	if botName == "" {
		botName = tenantCfg.DefaultBotName
	}
	avatarURL := req.AvatarURL
	if avatarURL == "" {
		avatarURL = tenantCfg.DefaultAvatarURL
	}
	lang := req.Language
	if lang == "" {
		lang = tenantCfg.DefaultLanguage
	}

	trxOpts := tenantCfg.TranscriptionOptions
	if req.TranscriptionOptions != nil {
		trxOpts = *req.TranscriptionOptions
	}

	autoLeave := tenantCfg.AutomaticLeave
	if req.AutomaticLeave != nil {
		autoLeave = *req.AutomaticLeave
	}

	// Cost Optimization Routing
	requiresPremium := req.VideoOptions != nil || req.EnableRealtimeAudioStream || req.Platform == PlatformWebex || req.Platform == PlatformMicrosoftTeams
	costProfile := s.costOptimizer.SelectOptimalProvider(core.ResourceTypeMeeting, requiresPremium)

	bot := &MeetingBotResponse{
		BotID:                   botID,
		TenantID:                tenantID,
		ProfileID:               req.ProfileID,
		MeetingURL:              req.MeetingURL,
		BotName:                 botName,
		AvatarURL:               avatarURL,
		Platform:                req.Platform,
		Status:                  StatusJoining,
		Language:                lang,
		RecordingMode:           req.RecordingMode,
		AudioStreamWebSocketURL: fmt.Sprintf("wss://api.patter.ai/ws/audio/%s", botID),
		TranscriptionOptions:    trxOpts,
		AutomaticLeave:          autoLeave,
		CreatedAt:               time.Now(),
		ProviderUsed:            costProfile.ProviderName,
	}

	s.mu.Lock()
	s.bots[tenantID+":"+botID] = bot
	s.mu.Unlock()

	if s.apiKey != "" {
		payload, _ := json.Marshal(map[string]interface{}{
			"meeting_url": req.MeetingURL,
			"bot_name":    botName,
			"avatar_url":  avatarURL,
			"transcription_options": map[string]interface{}{
				"provider": trxOpts.Provider,
				"language": lang,
			},
			"recording_mode": string(req.RecordingMode),
			"automatic_leave": map[string]interface{}{
				"everyone_left_timeout": autoLeave.EveryoneLeftTimeoutSec,
				"silence_timeout":       autoLeave.SilenceTimeoutSec,
			},
		})
		httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.patter.ai/v1/internal/bot", bytes.NewBuffer(payload))
		if err == nil {
			httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
			httpReq.Header.Set("Content-Type", "application/json")
			resp, err := s.httpClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				s.logger.Info("Patter meeting bot creation response", "status", resp.StatusCode, "len", len(body))
			}
		}
	}

	return bot, nil
}

func (s *meetingEngineService) GetBotStatus(ctx context.Context, tenantID, botID string) (*MeetingBotResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bot, ok := s.bots[tenantID+":"+botID]
	if !ok {
		return nil, fmt.Errorf("bot %s not found for tenant %s", botID, tenantID)
	}
	return bot, nil
}

func (s *meetingEngineService) ListBots(ctx context.Context, tenantID string) ([]*MeetingBotResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*MeetingBotResponse
	for _, bot := range s.bots {
		if bot.TenantID == tenantID {
			result = append(result, bot)
		}
	}
	return result, nil
}

func (s *meetingEngineService) SendChatMessage(ctx context.Context, tenantID string, req SendMeetingChatMessageRequest) error {
	s.mu.RLock()
	bot, ok := s.bots[tenantID+":"+req.BotID]
	s.mu.RUnlock()

	if !ok || bot.TenantID != tenantID {
		return fmt.Errorf("bot %s not found", req.BotID)
	}
	s.logger.Info("Sent chat message to meeting", "bot_id", req.BotID, "message", req.Message)
	return nil
}

func (s *meetingEngineService) LeaveMeeting(ctx context.Context, tenantID, botID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bot, ok := s.bots[tenantID+":"+botID]
	if !ok || bot.TenantID != tenantID {
		return fmt.Errorf("bot %s not found", botID)
	}
	bot.Status = StatusLeft
	return nil
}
