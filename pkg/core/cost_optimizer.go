package core

import (
	"log/slog"
)

type ResourceType string

const (
	ResourceTypeMeeting ResourceType = "MEETING"
	ResourceTypeVoice   ResourceType = "VOICE_CALL"
	ResourceTypeSTT     ResourceType = "STT"
	ResourceTypeTTS     ResourceType = "TTS"
)

type ProviderCostProfile struct {
	ProviderName string  `json:"provider_name"`
	CostPerMin   float64 `json:"cost_per_min"`
	CostPerReq   float64 `json:"cost_per_req"`
	Tier         string  `json:"tier"` // "ECONOMY_NATIVE", "STANDARD", "PREMIUM"
}

type CostOptimizationEngine struct {
	logger *slog.Logger
}

func NewCostOptimizationEngine(logger *slog.Logger) *CostOptimizationEngine {
	if logger == nil {
		logger = GetLogger()
	}
	return &CostOptimizationEngine{logger: logger.With("component", "cost_optimizer")}
}

func (e *CostOptimizationEngine) SelectOptimalProvider(resType ResourceType, requiresPremiumFeatures bool) ProviderCostProfile {
	switch resType {
	case ResourceTypeMeeting:
		if requiresPremiumFeatures {
			e.logger.Info("Selected Premium Meeting Provider for advanced video/audio features", "tier", "PREMIUM")
			return ProviderCostProfile{
				ProviderName: "patter_premium_engine",
				CostPerMin:   0.025,
				Tier:         "PREMIUM",
			}
		}
		e.logger.Info("Selected Economy Native Meeting Provider to minimize costs", "tier", "ECONOMY_NATIVE")
		return ProviderCostProfile{
			ProviderName: "patter_native_rtc",
			CostPerMin:   0.005,
			Tier:         "ECONOMY_NATIVE",
		}

	case ResourceTypeVoice:
		if requiresPremiumFeatures {
			return ProviderCostProfile{
				ProviderName: "patter_voice_hd_pro",
				CostPerMin:   0.03,
				Tier:         "PREMIUM",
			}
		}
		return ProviderCostProfile{
			ProviderName: "patter_sip_direct_pbx",
			CostPerMin:   0.002,
			Tier:         "ECONOMY_NATIVE",
		}

	case ResourceTypeTTS:
		if requiresPremiumFeatures {
			return ProviderCostProfile{
				ProviderName: "elevenlabs_hd",
				CostPerReq:   0.015,
				Tier:         "PREMIUM",
			}
		}
		return ProviderCostProfile{
			ProviderName: "openai_tts_standard",
			CostPerReq:   0.003,
			Tier:         "ECONOMY_NATIVE",
		}

	default:
		return ProviderCostProfile{
			ProviderName: "patter_native_default",
			CostPerMin:   0.001,
			Tier:         "ECONOMY_NATIVE",
		}
	}
}
