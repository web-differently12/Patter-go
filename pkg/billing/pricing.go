package billing

import (
	"sync"
	"time"
)

type ProviderCategory string

const (
	CategorySTT ProviderCategory = "STT"
	CategoryLLM ProviderCategory = "LLM"
	CategoryTTS ProviderCategory = "TTS"
	CategorySIP ProviderCategory = "SIP_TELECOM"
)

type ProviderUnit struct {
	Category       ProviderCategory `json:"category"`
	ProviderName   string           `json:"provider_name"`   // e.g., "deepgram", "openai", "elevenlabs", "twilio_sip"
	ModelName      string           `json:"model_name"`      // e.g., "nova-2", "gpt-4o", "flash-v2", "g711"
	RawCostPerUnit float64          `json:"raw_cost_per_unit"` // USD or EUR cost per ms/token/char/sec
	UnitType       string           `json:"unit_type"`       // "PER_MS", "PER_1K_INPUT_TOKENS", "PER_1K_OUTPUT_TOKENS", "PER_CHAR", "PER_SEC"
	OptionName     string           `json:"option_name,omitempty"` // e.g. "smart_format", "diarization", "voice_hd"
	OptionFee      float64          `json:"option_fee,omitempty"`
}

type TenantPricingConfig struct {
	TenantID         string    `json:"tenant_id"`
	MarkupPercentage float64   `json:"markup_percentage"` // e.g., 30.0 for +30%
	FixedFeePerCall  float64   `json:"fixed_fee_per_call"` // e.g., 0.02
	GraceLimitFiat   float64   `json:"grace_limit_fiat"`   // e.g., -0.50 EUR/USD allowed balance before graceful termination
	UpdatedAt        time.Time `json:"updated_at"`
}

type PricingGrid struct {
	mu            sync.RWMutex
	providers     map[string]ProviderUnit        // key: provider:model or provider:model:option
	tenantConfigs map[string]*TenantPricingConfig // key: tenant_id
}

func NewPricingGrid() *PricingGrid {
	pg := &PricingGrid{
		providers:     make(map[string]ProviderUnit),
		tenantConfigs: make(map[string]*TenantPricingConfig),
	}
	pg.loadDefaultProviderRates()
	return pg
}

func (pg *PricingGrid) loadDefaultProviderRates() {
	// Default Rates (USD/EUR)
	// STT
	pg.providers["deepgram:nova-2"] = ProviderUnit{
		Category:       CategorySTT,
		ProviderName:   "deepgram",
		ModelName:      "nova-2",
		RawCostPerUnit: 0.0000000043, // $0.0043 per minute => ~$0.000000071/sec
		UnitType:       "PER_MS",
	}
	pg.providers["assemblyai:conformer-2"] = ProviderUnit{
		Category:       CategorySTT,
		ProviderName:   "assemblyai",
		ModelName:      "conformer-2",
		RawCostPerUnit: 0.0000000062,
		UnitType:       "PER_MS",
	}
	// Add-on options
	pg.providers["deepgram:nova-2:diarization"] = ProviderUnit{
		Category:       CategorySTT,
		ProviderName:   "deepgram",
		ModelName:      "nova-2",
		RawCostPerUnit: 0.0000000010,
		UnitType:       "PER_MS",
		OptionName:     "diarization",
		OptionFee:      0.0000000010,
	}

	// LLM (per 1K tokens)
	pg.providers["openai:gpt-4o:input"] = ProviderUnit{
		Category:       CategoryLLM,
		ProviderName:   "openai",
		ModelName:      "gpt-4o",
		RawCostPerUnit: 0.0025, // $2.50 / 1M input tokens = $0.0025 / 1K
		UnitType:       "PER_1K_INPUT_TOKENS",
	}
	pg.providers["openai:gpt-4o:output"] = ProviderUnit{
		Category:       CategoryLLM,
		ProviderName:   "openai",
		ModelName:      "gpt-4o",
		RawCostPerUnit: 0.0100, // $10.00 / 1M output tokens = $0.01 / 1K
		UnitType:       "PER_1K_OUTPUT_TOKENS",
	}
	pg.providers["anthropic:claude-3-5-sonnet:input"] = ProviderUnit{
		Category:       CategoryLLM,
		ProviderName:   "anthropic",
		ModelName:      "claude-3-5-sonnet",
		RawCostPerUnit: 0.0030,
		UnitType:       "PER_1K_INPUT_TOKENS",
	}
	pg.providers["anthropic:claude-3-5-sonnet:output"] = ProviderUnit{
		Category:       CategoryLLM,
		ProviderName:   "anthropic",
		ModelName:      "claude-3-5-sonnet",
		RawCostPerUnit: 0.0150,
		UnitType:       "PER_1K_OUTPUT_TOKENS",
	}

	// TTS
	pg.providers["elevenlabs:flash-v2"] = ProviderUnit{
		Category:       CategoryTTS,
		ProviderName:   "elevenlabs",
		ModelName:      "flash-v2",
		RawCostPerUnit: 0.000015, // $15 per 1M chars
		UnitType:       "PER_CHAR",
	}
	pg.providers["cartesia:sonic"] = ProviderUnit{
		Category:       CategoryTTS,
		ProviderName:   "cartesia",
		ModelName:      "sonic",
		RawCostPerUnit: 0.000012,
		UnitType:       "PER_CHAR",
	}

	// SIP Telecom
	pg.providers["twilio:sip_trunk"] = ProviderUnit{
		Category:       CategorySIP,
		ProviderName:   "twilio",
		ModelName:      "sip_trunk",
		RawCostPerUnit: 0.00014, // ~$0.0085 per minute
		UnitType:       "PER_SEC",
	}
}

func (pg *PricingGrid) SetTenantConfig(config TenantPricingConfig) {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	config.UpdatedAt = time.Now()
	pg.tenantConfigs[config.TenantID] = &config
}

func (pg *PricingGrid) GetTenantConfig(tenantID string) TenantPricingConfig {
	pg.mu.RLock()
	defer pg.mu.RUnlock()

	if cfg, ok := pg.tenantConfigs[tenantID]; ok {
		return *cfg
	}
	return TenantPricingConfig{
		TenantID:         tenantID,
		MarkupPercentage: 25.0, // Default 25% margin
		FixedFeePerCall:  0.01,
		GraceLimitFiat:   -0.50,
		UpdatedAt:        time.Now(),
	}
}

func (pg *PricingGrid) GetProviderUnit(key string) (ProviderUnit, bool) {
	pg.mu.RLock()
	defer pg.mu.RUnlock()

	pu, ok := pg.providers[key]
	return pu, ok
}

func (pg *PricingGrid) GetAllProviderUnits() []ProviderUnit {
	pg.mu.RLock()
	defer pg.mu.RUnlock()

	units := make([]ProviderUnit, 0, len(pg.providers))
	for _, unit := range pg.providers {
		units = append(units, unit)
	}
	return units
}
