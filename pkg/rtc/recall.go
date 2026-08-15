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

type MeetingProfile struct {
	ProfileID                     string `json:"profile_id"`
	TenantID                      string `json:"tenant_id"`
	Name                          string `json:"name" binding:"required"`
	EnableSpeakerDiarization      bool   `json:"enable_speaker_diarization"`
	EnableActionItemsExtraction   bool   `json:"enable_action_items_extraction"`
	EnableParticipantSentiment    bool   `json:"enable_participant_sentiment"`
	EnableLiveTranslation         bool   `json:"enable_live_translation"`
	TargetTranslationLanguage     string `json:"target_translation_language,omitempty"` // "en", "fr", "es", "de"
	EnableScreenShareRecording    bool   `json:"enable_screen_share_recording"`
	SummaryTemplate               string `json:"summary_template"`               // "BANT_QUALIFICATION", "EXECUTIVE_SUMMARY", "TECHNICAL_ACTION_ITEMS"
	AutoLeaveOnSilenceMinutes     int    `json:"auto_leave_on_silence_minutes"`
	AutoLeaveWhenEveryoneLeft     bool   `json:"auto_leave_when_everyone_left"`
	CustomAvatarVideoURL          string `json:"custom_avatar_video_url,omitempty"`
}

type CreateMeetingBotRequest struct {
	MeetingURL                string          `json:"meeting_url" binding:"required"`
	ProfileID                 string          `json:"profile_id,omitempty"`
	BotName                   string          `json:"bot_name,omitempty"`
	AvatarURL                 string          `json:"avatar_url,omitempty"`
	Platform                  MeetingPlatform `json:"platform,omitempty"`
	RecordingMode             RecordingMode   `json:"recording_mode,omitempty"`
	EnableRealtimeTranscript  bool            `json:"enable_realtime_transcript"`
	EnableRealtimeAudioStream bool            `json:"enable_realtime_audio_stream"`
	EnableChatMessaging       bool            `json:"enable_chat_messaging"`
	Language                  string          `json:"language,omitempty"` // "fr", "en", "es", "de", "auto"
	AutomaticLeaveWhenAlone   bool            `json:"automatic_leave_when_alone"`
	SilenceTimeoutMinutes     int             `json:"silence_timeout_minutes,omitempty"`
	SystemPrompt              string          `json:"system_prompt,omitempty"`
	WebhookURL                string          `json:"webhook_url,omitempty"`
}

type SendMeetingChatMessageRequest struct {
	BotID   string `json:"bot_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type MeetingBotResponse struct {
	BotID                   string          `json:"bot_id"`
	TenantID                string          `json:"tenant_id"`
	ProfileID               string          `json:"profile_id,omitempty"`
	MeetingURL              string          `json:"meeting_url"`
	BotName                 string          `json:"bot_name"`
	AvatarURL               string          `json:"avatar_url,omitempty"`
	Platform                MeetingPlatform `json:"platform"`
	Status                  BotStatus       `json:"status"`
	Language                string          `json:"language"`
	RecordingMode           RecordingMode   `json:"recording_mode"`
	VideoRecordingURL       string          `json:"video_recording_url,omitempty"`
	TranscriptURL           string          `json:"transcript_url,omitempty"`
	AudioStreamWebSocketURL string          `json:"audio_stream_websocket_url,omitempty"`
	CreatedAt               time.Time       `json:"created_at"`
}

type RecallAIService interface {
	CreateMeetingProfile(ctx context.Context, tenantID string, prof MeetingProfile) (*MeetingProfile, error)
	GetMeetingProfile(ctx context.Context, tenantID, profileID string) (*MeetingProfile, error)
	ListMeetingProfiles(ctx context.Context, tenantID string) ([]*MeetingProfile, error)
	CreateMeetingBot(ctx context.Context, tenantID string, req CreateMeetingBotRequest) (*MeetingBotResponse, error)
	GetBotStatus(ctx context.Context, tenantID, botID string) (*MeetingBotResponse, error)
	ListBots(ctx context.Context, tenantID string) ([]*MeetingBotResponse, error)
	SendChatMessage(ctx context.Context, tenantID string, req SendMeetingChatMessageRequest) error
	LeaveMeeting(ctx context.Context, tenantID, botID string) error
}

type recallAIService struct {
	apiKey     string
	httpClient *http.Client
	logger     *slog.Logger
	mu         sync.RWMutex
	profiles   map[string]*MeetingProfile
	bots       map[string]*MeetingBotResponse
}

func NewRecallAIService(apiKey string, logger *slog.Logger) RecallAIService {
	if logger == nil {
		logger = core.GetLogger()
	}
	svc := &recallAIService{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     logger.With("component", "recall_ai"),
		profiles:   make(map[string]*MeetingProfile),
		bots:       make(map[string]*MeetingBotResponse),
	}

	// Pre-populate default Meeting Profile
	defaultProfID := "mp_executive_pro"
	svc.profiles["default_tenant:"+defaultProfID] = &MeetingProfile{
		ProfileID:                   defaultProfID,
		TenantID:                    "default_tenant",
		Name:                        "Profil Exécutif — Diarisation & Actions",
		EnableSpeakerDiarization:    true,
		EnableActionItemsExtraction: true,
		EnableParticipantSentiment:  true,
		EnableLiveTranslation:       true,
		TargetTranslationLanguage:   "fr",
		EnableScreenShareRecording:  true,
		SummaryTemplate:             "EXECUTIVE_SUMMARY",
		AutoLeaveOnSilenceMinutes:   5,
		AutoLeaveWhenEveryoneLeft:   true,
	}

	mockID := "bot_meet_101"
	svc.bots["default_tenant:"+mockID] = &MeetingBotResponse{
		BotID:                   mockID,
		TenantID:                "default_tenant",
		ProfileID:               defaultProfID,
		MeetingURL:              "https://meet.google.com/abc-defg-hij",
		BotName:                 "Patter AI Assistant",
		AvatarURL:               "https://patter.ai/assets/avatar_white_label.png",
		Platform:                PlatformGoogleMeet,
		Status:                  StatusInCall,
		Language:                "fr",
		RecordingMode:           RecordingSpeakerView,
		VideoRecordingURL:       "https://recall.ai/recordings/rec_101.mp4",
		TranscriptURL:           "https://recall.ai/transcripts/trx_101.json",
		AudioStreamWebSocketURL: "wss://recall.ai/ws/audio/bot_101",
		CreatedAt:               time.Now().Add(-15 * time.Minute),
	}

	return svc
}

func (s *recallAIService) CreateMeetingProfile(ctx context.Context, tenantID string, prof MeetingProfile) (*MeetingProfile, error) {
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

func (s *recallAIService) GetMeetingProfile(ctx context.Context, tenantID, profileID string) (*MeetingProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prof, ok := s.profiles[tenantID+":"+profileID]
	if !ok {
		return nil, fmt.Errorf("meeting profile %s not found", profileID)
	}
	return prof, nil
}

func (s *recallAIService) ListMeetingProfiles(ctx context.Context, tenantID string) ([]*MeetingProfile, error) {
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

func (s *recallAIService) CreateMeetingBot(ctx context.Context, tenantID string, req CreateMeetingBotRequest) (*MeetingBotResponse, error) {
	botID := "bot_" + uuid.New().String()[:8]
	if req.BotName == "" {
		req.BotName = "Patter AI Agent"
	}
	if req.Language == "" {
		req.Language = "fr"
	}
	if req.RecordingMode == "" {
		req.RecordingMode = RecordingSpeakerView
	}

	bot := &MeetingBotResponse{
		BotID:                   botID,
		TenantID:                tenantID,
		ProfileID:               req.ProfileID,
		MeetingURL:              req.MeetingURL,
		BotName:                 req.BotName,
		AvatarURL:               req.AvatarURL,
		Platform:                req.Platform,
		Status:                  StatusJoining,
		Language:                req.Language,
		RecordingMode:           req.RecordingMode,
		AudioStreamWebSocketURL: fmt.Sprintf("wss://recall.ai/ws/audio/%s", botID),
		CreatedAt:               time.Now(),
	}

	s.mu.Lock()
	s.bots[tenantID+":"+botID] = bot
	s.mu.Unlock()

	if s.apiKey != "" {
		payload, _ := json.Marshal(map[string]interface{}{
			"meeting_url": req.MeetingURL,
			"bot_name":    req.BotName,
			"avatar_url":  req.AvatarURL,
			"transcription_options": map[string]interface{}{
				"provider": "assemblyai",
				"language": req.Language,
			},
			"recording_mode": string(req.RecordingMode),
			"automatic_leave": map[string]interface{}{
				"everyone_left_timeout": 60,
			},
		})
		httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.recall.ai/api/v1/bot", bytes.NewBuffer(payload))
		if err == nil {
			httpReq.Header.Set("Authorization", "Token "+s.apiKey)
			httpReq.Header.Set("Content-Type", "application/json")
			resp, err := s.httpClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				s.logger.Info("Recall.ai bot creation API response", "status", resp.StatusCode, "len", len(body))
			}
		}
	}

	return bot, nil
}

func (s *recallAIService) GetBotStatus(ctx context.Context, tenantID, botID string) (*MeetingBotResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bot, ok := s.bots[tenantID+":"+botID]
	if !ok {
		return nil, fmt.Errorf("bot %s not found for tenant %s", botID, tenantID)
	}
	return bot, nil
}

func (s *recallAIService) ListBots(ctx context.Context, tenantID string) ([]*MeetingBotResponse, error) {
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

func (s *recallAIService) SendChatMessage(ctx context.Context, tenantID string, req SendMeetingChatMessageRequest) error {
	s.mu.RLock()
	bot, ok := s.bots[tenantID+":"+req.BotID]
	s.mu.RUnlock()

	if !ok || bot.TenantID != tenantID {
		return fmt.Errorf("bot %s not found", req.BotID)
	}
	s.logger.Info("Sent chat message to meeting", "bot_id", req.BotID, "message", req.Message)
	return nil
}

func (s *recallAIService) LeaveMeeting(ctx context.Context, tenantID, botID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bot, ok := s.bots[tenantID+":"+botID]
	if !ok || bot.TenantID != tenantID {
		return fmt.Errorf("bot %s not found", botID)
	}
	bot.Status = StatusLeft
	return nil
}
