package billing

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/core"
)

type EventChannel string

const (
	ChannelVoiceCall  EventChannel = "voice_call"
	ChannelMeetingBot EventChannel = "meeting_bot"
	ChannelWhatsApp   EventChannel = "whatsapp"
	ChannelMCPTool    EventChannel = "mcp_tool"
)

type MeteringEvent struct {
	EventID          string       `json:"event_id"`
	TenantID         string       `json:"tenant_id"`
	AgentID          string       `json:"agent_id,omitempty"`
	SessionID        string       `json:"session_id"`
	ChannelType      EventChannel `json:"channel_type"`
	Provider         string       `json:"provider"`   // e.g. "deepgram", "openai", "elevenlabs", "twilio"
	ModelName        string       `json:"model_name"` // e.g. "nova-2", "gpt-4o", "flash-v2", "sip_trunk"
	PromptTokens     int          `json:"prompt_tokens,omitempty"`
	CompletionTokens int          `json:"completion_tokens,omitempty"`
	AudioDurationMs  int          `json:"audio_duration_ms,omitempty"`
	TTSCharacters    int          `json:"tts_characters,omitempty"`
	SIPDurationSec   int          `json:"sip_duration_sec,omitempty"`
	AddonOptions     []string     `json:"addon_options,omitempty"` // e.g. ["diarization", "smart_format"]
	Timestamp        time.Time    `json:"timestamp"`
}

type MeteringResult struct {
	EventID                 string    `json:"event_id"`
	TenantID                string    `json:"tenant_id"`
	SessionID               string    `json:"session_id"`
	RawCostUSD              float64   `json:"raw_cost_usd"`
	BilledCostUSD           float64   `json:"billed_cost_usd"`
	ProfitUSD               float64   `json:"profit_usd"`
	WalletBalanceAfter      float64   `json:"wallet_balance_after"`
	SessionGracefulShutdown bool      `json:"session_graceful_shutdown"` // Signal SESSION_TERMINATION_GRACEFUL
	Timestamp               time.Time `json:"timestamp"`
}

type UsageMeterService interface {
	IngestEvent(ctx context.Context, evt MeteringEvent) (*MeteringResult, error)
	GetLiveSessionMetrics(sessionID string) (*MeteringResult, bool)
}

type usageMeterService struct {
	mu           sync.RWMutex
	pricingGrid  *PricingGrid
	walletLedger WalletLedgerService
	liveMetrics  map[string]*MeteringResult
	logger       *slog.Logger
}

func NewUsageMeterService(pricingGrid *PricingGrid, walletLedger WalletLedgerService, logger *slog.Logger) UsageMeterService {
	if logger == nil {
		logger = core.GetLogger()
	}
	return &usageMeterService{
		pricingGrid:  pricingGrid,
		walletLedger: walletLedger,
		liveMetrics:  make(map[string]*MeteringResult),
		logger:       logger.With("component", "usage_meter_engine"),
	}
}

func (s *usageMeterService) IngestEvent(ctx context.Context, evt MeteringEvent) (*MeteringResult, error) {
	startTime := time.Now()
	if evt.EventID == "" {
		evt.EventID = "evt_" + uuid.New().String()[:8]
	}
	if evt.Timestamp.IsZero() {
		evt.Timestamp = time.Now()
	}

	tenantCfg := s.pricingGrid.GetTenantConfig(evt.TenantID)

	var rawCost float64

	// 1. STT Audio Duration
	if evt.AudioDurationMs > 0 {
		key := fmt.Sprintf("%s:%s", evt.Provider, evt.ModelName)
		if pu, ok := s.pricingGrid.GetProviderUnit(key); ok {
			rawCost += float64(evt.AudioDurationMs) * pu.RawCostPerUnit
		} else {
			rawCost += float64(evt.AudioDurationMs) * 0.0000000043 // default Nova-2 fallback
		}
	}

	// 2. LLM Prompt & Completion Tokens
	if evt.PromptTokens > 0 {
		key := fmt.Sprintf("%s:%s:input", evt.Provider, evt.ModelName)
		if pu, ok := s.pricingGrid.GetProviderUnit(key); ok {
			rawCost += (float64(evt.PromptTokens) / 1000.0) * pu.RawCostPerUnit
		} else {
			rawCost += (float64(evt.PromptTokens) / 1000.0) * 0.0025 // default GPT-4o input
		}
	}
	if evt.CompletionTokens > 0 {
		key := fmt.Sprintf("%s:%s:output", evt.Provider, evt.ModelName)
		if pu, ok := s.pricingGrid.GetProviderUnit(key); ok {
			rawCost += (float64(evt.CompletionTokens) / 1000.0) * pu.RawCostPerUnit
		} else {
			rawCost += (float64(evt.CompletionTokens) / 1000.0) * 0.0100 // default GPT-4o output
		}
	}

	// 3. TTS Characters
	if evt.TTSCharacters > 0 {
		key := fmt.Sprintf("%s:%s", evt.Provider, evt.ModelName)
		if pu, ok := s.pricingGrid.GetProviderUnit(key); ok {
			rawCost += float64(evt.TTSCharacters) * pu.RawCostPerUnit
		} else {
			rawCost += float64(evt.TTSCharacters) * 0.000015 // default ElevenLabs
		}
	}

	// 4. SIP Telecom Seconds
	if evt.SIPDurationSec > 0 {
		key := fmt.Sprintf("%s:%s", evt.Provider, evt.ModelName)
		if pu, ok := s.pricingGrid.GetProviderUnit(key); ok {
			rawCost += float64(evt.SIPDurationSec) * pu.RawCostPerUnit
		} else {
			rawCost += float64(evt.SIPDurationSec) * 0.00014 // default Twilio SIP
		}
	}

	// 5. Addon options fee
	for _, opt := range evt.AddonOptions {
		key := fmt.Sprintf("%s:%s:%s", evt.Provider, evt.ModelName, opt)
		if pu, ok := s.pricingGrid.GetProviderUnit(key); ok {
			rawCost += pu.OptionFee
		}
	}

	// Tenant Billing Calculation Formula:
	// Billed = (RawCost * (1 + Markup% / 100)) + FixedFee
	billedCost := (rawCost * (1.0 + (tenantCfg.MarkupPercentage / 100.0)))
	profit := billedCost - rawCost

	// Deduct from Wallet Ledger
	balAfter, err := s.walletLedger.DeductUsage(ctx, evt.TenantID, evt.SessionID, billedCost, evt)
	if err != nil {
		s.logger.Error("Wallet ledger deduction failed", "error", err, "tenant_id", evt.TenantID)
	}

	// Check Grace Limit Kill-Switch
	gracefulShutdown := false
	if balAfter <= tenantCfg.GraceLimitFiat {
		gracefulShutdown = true
		s.logger.Warn("SESSION_TERMINATION_GRACEFUL triggered due to exhausted tenant wallet",
			"tenant_id", evt.TenantID, "session_id", evt.SessionID, "balance_after", balAfter, "grace_limit", tenantCfg.GraceLimitFiat)
	}

	result := &MeteringResult{
		EventID:                 evt.EventID,
		TenantID:                evt.TenantID,
		SessionID:               evt.SessionID,
		RawCostUSD:              rawCost,
		BilledCostUSD:           billedCost,
		ProfitUSD:               profit,
		WalletBalanceAfter:      balAfter,
		SessionGracefulShutdown: gracefulShutdown,
		Timestamp:               time.Now(),
	}

	s.mu.Lock()
	if existing, ok := s.liveMetrics[evt.SessionID]; ok {
		existing.RawCostUSD += result.RawCostUSD
		existing.BilledCostUSD += result.BilledCostUSD
		existing.ProfitUSD += result.ProfitUSD
		existing.WalletBalanceAfter = result.WalletBalanceAfter
		existing.SessionGracefulShutdown = result.SessionGracefulShutdown
		existing.Timestamp = result.Timestamp
	} else {
		s.liveMetrics[evt.SessionID] = &MeteringResult{
			EventID:                 result.EventID,
			TenantID:                result.TenantID,
			SessionID:               result.SessionID,
			RawCostUSD:              result.RawCostUSD,
			BilledCostUSD:           result.BilledCostUSD,
			ProfitUSD:               result.ProfitUSD,
			WalletBalanceAfter:      result.WalletBalanceAfter,
			SessionGracefulShutdown: result.SessionGracefulShutdown,
			Timestamp:               result.Timestamp,
		}
	}
	s.mu.Unlock()

	elapsed := time.Since(startTime)
	s.logger.Debug("Metering calculation completed", "event_id", evt.EventID, "elapsed_us", elapsed.Microseconds(), "billed_usd", billedCost)

	return result, nil
}

func (s *usageMeterService) GetLiveSessionMetrics(sessionID string) (*MeteringResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res, ok := s.liveMetrics[sessionID]
	if !ok {
		return nil, false
	}
	return res, true
}
