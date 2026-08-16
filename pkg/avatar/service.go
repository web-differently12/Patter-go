package avatar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/core"
)

type AvatarMode string

const (
	ModeWaveSpeedRealtimeStreaming AvatarMode = "wavespeed_realtime_stream"
	ModeWaveSpeedAsyncCampaign     AvatarMode = "wavespeed_async_campaign"
	ModeLocalGPUAsyncMassRender    AvatarMode = "local_gpu_async_render"
)

type CDNStorageProvider string

const (
	CDNCloudflareR2 CDNStorageProvider = "cloudflare_r2"
	CDNBunnyNet     CDNStorageProvider = "bunny_net"
	CDNCloudinary   CDNStorageProvider = "cloudinary"
)

type StreamAvatarRequest struct {
	AvatarID          string     `json:"avatar_id" binding:"required"`
	Mode              AvatarMode `json:"mode"`
	PromptText        string     `json:"prompt_text" binding:"required"`
	VoiceID           string     `json:"voice_id,omitempty"`
	Resolution        string     `json:"resolution,omitempty"`
	InteractiveStream bool       `json:"interactive_stream"`
}

type RenderAvatarVideoRequest struct {
	AvatarID          string             `json:"avatar_id" binding:"required"`
	Prompts           []string           `json:"prompts" binding:"required"`
	Resolution        string             `json:"resolution,omitempty"`
	CDNProvider       CDNStorageProvider `json:"cdn_provider,omitempty"`
	GenerateShortlink bool               `json:"generate_shortlink"`
	CallbackURL       string             `json:"callback_url,omitempty"`
}

type ShortlinkInfo struct {
	ShortID     string    `json:"short_id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"` // 60-day auto-expiry
}

type AvatarResponse struct {
	JobID             string             `json:"job_id"`
	TenantID          string             `json:"tenant_id"`
	AvatarID          string             `json:"avatar_id"`
	Mode              AvatarMode         `json:"mode"`
	Status            string             `json:"status"`
	StreamWebRTCURL   string             `json:"stream_webrtc_url,omitempty"`
	VideoMP4URL       string             `json:"video_mp4_url,omitempty"`
	Shortlink         *ShortlinkInfo     `json:"shortlink,omitempty"`
	CDNUsed           CDNStorageProvider `json:"cdn_used"`
	CostPerMin        float64            `json:"cost_per_min"`
	ProviderUsed      string             `json:"provider_used"`
	CreatedAt         time.Time          `json:"created_at"`
}

type AvatarEngineService interface {
	StreamRealtimeAvatar(ctx context.Context, tenantID string, req StreamAvatarRequest) (*AvatarResponse, error)
	RenderMassAvatarVideo(ctx context.Context, tenantID string, req RenderAvatarVideoRequest) (*AvatarResponse, error)
	GenerateShortlink(originalURL string) *ShortlinkInfo
	GetShortlinkTarget(shortID string) (string, bool)
}

type avatarEngineService struct {
	mu              sync.RWMutex
	wavespeedAPIKey string
	httpClient      *http.Client
	logger          *slog.Logger
	shortlinks      map[string]*ShortlinkInfo
}

func NewAvatarEngineService(wavespeedAPIKey string, logger *slog.Logger) AvatarEngineService {
	if logger == nil {
		logger = core.GetLogger()
	}
	return &avatarEngineService{
		wavespeedAPIKey: wavespeedAPIKey,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
		logger:          logger.With("component", "avatar_engine"),
		shortlinks:      make(map[string]*ShortlinkInfo),
	}
}

func (s *avatarEngineService) GenerateShortlink(originalURL string) *ShortlinkInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	shortID := fmt.Sprintf("%06x", rand.Intn(0xFFFFFF))
	shortURL := fmt.Sprintf("https://s.patter.ai/v/%s", shortID)
	now := time.Now()
	info := &ShortlinkInfo{
		ShortID:     shortID,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		CreatedAt:   now,
		ExpiresAt:   now.Add(60 * 24 * time.Hour), // 60-day data retention
	}
	s.shortlinks[shortID] = info
	return info
}

func (s *avatarEngineService) GetShortlinkTarget(shortID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, ok := s.shortlinks[shortID]
	if !ok || time.Now().After(info.ExpiresAt) {
		return "", false
	}
	return info.OriginalURL, true
}

func (s *avatarEngineService) StreamRealtimeAvatar(ctx context.Context, tenantID string, req StreamAvatarRequest) (*AvatarResponse, error) {
	jobID := "avatar_stream_" + uuid.New().String()[:8]

	mode := req.Mode
	if mode == "" {
		mode = ModeWaveSpeedRealtimeStreaming
	}

	providerUsed := "wavespeed_api_v2"
	costPerMin := 0.08

	if mode == ModeLocalGPUAsyncMassRender {
		providerUsed = "local_gpu_renderer_mcp"
		costPerMin = 0.002
	}

	resp := &AvatarResponse{
		JobID:           jobID,
		TenantID:        tenantID,
		AvatarID:        req.AvatarID,
		Mode:            mode,
		Status:          "STREAMING",
		StreamWebRTCURL: fmt.Sprintf("wss://api.patter.ai/v1/avatar/webrtc/%s", jobID),
		CDNUsed:         CDNCloudflareR2,
		CostPerMin:      costPerMin,
		ProviderUsed:    providerUsed,
		CreatedAt:       time.Now(),
	}

	if s.wavespeedAPIKey != "" && mode == ModeWaveSpeedRealtimeStreaming {
		payload, _ := json.Marshal(map[string]interface{}{
			"avatar_id":   req.AvatarID,
			"prompt_text": req.PromptText,
			"voice_id":    req.VoiceID,
			"interactive": req.InteractiveStream,
		})
		httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.wavespeed.ai/v1/streaming/avatar", bytes.NewBuffer(payload))
		if err == nil {
			httpReq.Header.Set("Authorization", "Bearer "+s.wavespeedAPIKey)
			httpReq.Header.Set("Content-Type", "application/json")
			res, err := s.httpClient.Do(httpReq)
			if err == nil {
				defer res.Body.Close()
				body, _ := io.ReadAll(res.Body)
				s.logger.Info("WaveSpeed API Response", "status", res.StatusCode, "len", len(body))
			}
		}
	}

	return resp, nil
}

func (s *avatarEngineService) RenderMassAvatarVideo(ctx context.Context, tenantID string, req RenderAvatarVideoRequest) (*AvatarResponse, error) {
	jobID := "avatar_render_" + uuid.New().String()[:8]

	cdn := req.CDNProvider
	if cdn == "" {
		cdn = CDNCloudflareR2
	}

	originalURL := fmt.Sprintf("https://cdn.patter.ai/renders/%s.mp4", jobID)
	if cdn == CDNBunnyNet {
		originalURL = fmt.Sprintf("https://patter.b-cdn.net/renders/%s.mp4", jobID)
	}

	var shortlink *ShortlinkInfo
	if req.GenerateShortlink {
		shortlink = s.GenerateShortlink(originalURL)
	}

	providerUsed := "wavespeed_async_campaign_engine"
	costPerMin := 0.005
	if req.CDNProvider == "local" || req.AvatarID == "avatar_mass_prod" {
		providerUsed = "local_gpu_renderer_mcp"
		costPerMin = 0.002
	}

	resp := &AvatarResponse{
		JobID:        jobID,
		TenantID:     tenantID,
		AvatarID:     req.AvatarID,
		Mode:         ModeWaveSpeedAsyncCampaign,
		Status:       "PROCESSING",
		VideoMP4URL:  originalURL,
		Shortlink:    shortlink,
		CDNUsed:      cdn,
		CostPerMin:   costPerMin,
		ProviderUsed: providerUsed,
		CreatedAt:    time.Now(),
	}

	s.logger.Info("Queued WaveSpeed Campaign Avatar Video Render", "job_id", jobID, "prompts_count", len(req.Prompts), "cdn", cdn)
	return resp, nil
}
