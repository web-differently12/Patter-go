package rtc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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

type BotStatus string

const (
	StatusJoining   BotStatus = "joining"
	StatusInCall    BotStatus = "in_call"
	StatusRecording BotStatus = "recording"
	StatusLeft      BotStatus = "left"
	StatusFailed    BotStatus = "failed"
)

type CreateMeetingBotRequest struct {
	MeetingURL          string          `json:"meeting_url" binding:"required"`
	BotName             string          `json:"bot_name,omitempty"`
	Platform            MeetingPlatform `json:"platform,omitempty"`
	EnableRecording     bool            `json:"enable_recording"`
	EnableLiveStreaming bool            `json:"enable_live_streaming"`
	Language            string          `json:"language,omitempty"` // e.g. "fr", "en", "es", "auto"
	SystemPrompt        string          `json:"system_prompt,omitempty"`
	WebhookURL          string          `json:"webhook_url,omitempty"`
}

type MeetingBotResponse struct {
	BotID               string          `json:"bot_id"`
	TenantID            string          `json:"tenant_id"`
	MeetingURL          string          `json:"meeting_url"`
	BotName             string          `json:"bot_name"`
	Platform            MeetingPlatform `json:"platform"`
	Status              BotStatus       `json:"status"`
	Language            string          `json:"language"`
	VideoRecordingURL   string          `json:"video_recording_url,omitempty"`
	TranscriptURL       string          `json:"transcript_url,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
}

type RecallAIService interface {
	CreateMeetingBot(ctx context.Context, tenantID string, req CreateMeetingBotRequest) (*MeetingBotResponse, error)
	GetBotStatus(ctx context.Context, tenantID, botID string) (*MeetingBotResponse, error)
	ListBots(ctx context.Context, tenantID string) ([]*MeetingBotResponse, error)
	LeaveMeeting(ctx context.Context, tenantID, botID string) error
}

type recallAIService struct {
	apiKey     string
	httpClient *http.Client
	logger     *slog.Logger
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
		bots:       make(map[string]*MeetingBotResponse),
	}

	// Pre-populate mock meeting bot for demonstration
	mockID := "bot_meet_101"
	svc.bots["default_tenant:"+mockID] = &MeetingBotResponse{
		BotID:             mockID,
		TenantID:          "default_tenant",
		MeetingURL:        "https://meet.google.com/abc-defg-hij",
		BotName:           "Patter AI Assistant",
		Platform:          PlatformGoogleMeet,
		Status:            StatusInCall,
		Language:          "fr",
		VideoRecordingURL: "https://recall.ai/recordings/rec_101.mp4",
		TranscriptURL:     "https://recall.ai/transcripts/trx_101.json",
		CreatedAt:         time.Now().Add(-15 * time.Minute),
	}

	return svc
}

func (s *recallAIService) CreateMeetingBot(ctx context.Context, tenantID string, req CreateMeetingBotRequest) (*MeetingBotResponse, error) {
	botID := "bot_" + uuid.New().String()[:8]
	if req.BotName == "" {
		req.BotName = "Patter AI Agent"
	}
	if req.Language == "" {
		req.Language = "fr"
	}

	bot := &MeetingBotResponse{
		BotID:             botID,
		TenantID:          tenantID,
		MeetingURL:        req.MeetingURL,
		BotName:           req.BotName,
		Platform:          req.Platform,
		Status:            StatusJoining,
		Language:          req.Language,
		CreatedAt:         time.Now(),
	}

	s.bots[tenantID+":"+botID] = bot

	// Recall.ai API integration
	if s.apiKey != "" {
		payload, _ := json.Marshal(map[string]interface{}{
			"meeting_url": req.MeetingURL,
			"bot_name":    req.BotName,
			"transcription_options": map[string]interface{}{
				"provider": "assemblyai",
				"language": req.Language,
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
	bot, ok := s.bots[tenantID+":"+botID]
	if !ok {
		return nil, fmt.Errorf("bot %s not found for tenant %s", botID, tenantID)
	}
	return bot, nil
}

func (s *recallAIService) ListBots(ctx context.Context, tenantID string) ([]*MeetingBotResponse, error) {
	var result []*MeetingBotResponse
	for _, bot := range s.bots {
		if bot.TenantID == tenantID {
			result = append(result, bot)
		}
	}
	return result, nil
}

func (s *recallAIService) LeaveMeeting(ctx context.Context, tenantID, botID string) error {
	bot, ok := s.bots[tenantID+":"+botID]
	if !ok {
		return fmt.Errorf("bot %s not found", botID)
	}
	bot.Status = StatusLeft
	return nil
}
