package core

import (
	"log/slog"
	"strings"
	"time"
)

type ResourceType string

const (
	ResourceTypeMeeting ResourceType = "MEETING"
	ResourceTypeVoice   ResourceType = "VOICE_CALL"
	ResourceTypeSTT     ResourceType = "STT"
	ResourceTypeTTS     ResourceType = "TTS"
	ResourceTypeAvatar  ResourceType = "AVATAR"
	ResourceTypeLLM     ResourceType = "LLM"
)

type ProviderCostProfile struct {
	ProviderName string            `json:"provider_name"`
	CostPerMin   float64           `json:"cost_per_min"`
	CostPerReq   float64           `json:"cost_per_req"`
	Tier         string            `json:"tier"` // "ECONOMY_NATIVE", "STANDARD", "PREMIUM"
	FeatureFlags map[string]bool   `json:"feature_flags"`
	FallbackChain []string         `json:"fallback_chain"`
}

type ProviderMatrixEntry struct {
	ProviderName string          `json:"provider_name"`
	ResourceType ResourceType    `json:"resource_type"`
	CostPerMin   float64         `json:"cost_per_min"`
	CostPerReq   float64         `json:"cost_per_req"`
	OffPeakDiscount float64      `json:"off_peak_discount"` // e.g. 0.5 for -50% off-peak
	FeatureFlags map[string]bool `json:"feature_flags"`
}

type CostOptimizationEngine struct {
	logger   *slog.Logger
	matrix   []ProviderMatrixEntry
}

func NewCostOptimizationEngine(logger *slog.Logger) *CostOptimizationEngine {
	if logger == nil {
		logger = GetLogger()
	}

	// Matrix table with real-time costs and feature flags
	matrix := []ProviderMatrixEntry{
		{
			ProviderName: "patter_native_rtc",
			ResourceType: ResourceTypeMeeting,
			CostPerMin:   0.005,
			FeatureFlags: map[string]bool{"audio_recording": true, "transcription": true},
		},
		{
			ProviderName: "meetingbaas_engine",
			ResourceType: ResourceTypeMeeting,
			CostPerMin:   0.010,
			FeatureFlags: map[string]bool{"audio_recording": true, "transcription": true, "google_meet": true, "zoom": true},
		},
		{
			ProviderName: "patter_premium_engine",
			ResourceType: ResourceTypeMeeting,
			CostPerMin:   0.025,
			FeatureFlags: map[string]bool{"audio_recording": true, "video_hd": true, "transcription": true, "realtime_stream": true, "diarization": true, "translation": true, "webex": true, "teams": true},
		},
		{
			ProviderName: "patter_sip_direct_pbx",
			ResourceType: ResourceTypeVoice,
			CostPerMin:   0.002,
			FeatureFlags: map[string]bool{"sip_trunk": true, "raw_rtp": true},
		},
		{
			ProviderName: "patter_voice_hd_pro",
			ResourceType: ResourceTypeVoice,
			CostPerMin:   0.03,
			FeatureFlags: map[string]bool{"sip_trunk": true, "barge_in": true, "krisp_noise_suppression": true, "ambient_audio": true},
		},
		{
			ProviderName: "local_gpu_renderer_mcp",
			ResourceType: ResourceTypeAvatar,
			CostPerMin:   0.002,
			FeatureFlags: map[string]bool{"async_render": true, "mass_generation": true},
		},
		{
			ProviderName: "wavespeed_api_v2",
			ResourceType: ResourceTypeAvatar,
			CostPerMin:   0.08,
			FeatureFlags: map[string]bool{"realtime_webrtc": true, "interactive_avatar": true, "hd_stream": true},
		},
		{
			ProviderName: "deepseek_r1_standard",
			ResourceType: ResourceTypeLLM,
			CostPerReq:   0.001,
			OffPeakDiscount: 0.50, // -50% discount during off-peak hours
			FeatureFlags: map[string]bool{"reasoning": true, "off_peak_discount": true},
		},
		{
			ProviderName: "qwen_2_5_turbo",
			ResourceType: ResourceTypeLLM,
			CostPerReq:   0.002,
			FeatureFlags: map[string]bool{"multilingual": true, "fast_response": true},
		},
	}

	return &CostOptimizationEngine{
		logger: logger.With("component", "cost_optimizer"),
		matrix: matrix,
	}
}

func (e *CostOptimizationEngine) ResolveProvider(resType ResourceType, requiredFeatures map[string]bool) ProviderCostProfile {
	now := time.Now()
	isOffPeak := now.Hour() >= 22 || now.Hour() <= 6

	var bestEntry *ProviderMatrixEntry
	lowestCost := 9999.0

	for _, entry := range e.matrix {
		if entry.ResourceType != resType {
			continue
		}

		// Check if entry satisfies all required features
		hasAllFeatures := true
		for reqFeat, required := range requiredFeatures {
			if required && !entry.FeatureFlags[reqFeat] {
				hasAllFeatures = false;
				break
			}
		}

		if !hasAllFeatures {
			continue
		}

		effectiveCost := entry.CostPerMin + entry.CostPerReq
		if isOffPeak && entry.OffPeakDiscount > 0 {
			effectiveCost = effectiveCost * (1.0 - entry.OffPeakDiscount)
		}

		if effectiveCost < lowestCost {
			lowestCost = effectiveCost
			bestEntry = &entry
		}
	}

	if bestEntry == nil {
		e.logger.Info("No custom match found in provider matrix, selecting default fallback", "resource", string(resType))
		return ProviderCostProfile{
			ProviderName: "patter_default_fallback",
			CostPerMin:   0.01,
			Tier:         "ECONOMY_NATIVE",
			FallbackChain: []string{"patter_default_fallback", "patter_secondary_fallback"},
		}
	}

	tier := "ECONOMY_NATIVE"
	if strings.Contains(bestEntry.ProviderName, "premium") || strings.Contains(bestEntry.ProviderName, "wavespeed") {
		tier = "PREMIUM"
	} else if strings.Contains(bestEntry.ProviderName, "baas") || strings.Contains(bestEntry.ProviderName, "hd") {
		tier = "STANDARD"
	}

	e.logger.Info("Resolved lowest cost provider from matrix",
		"resource", string(resType),
		"selected_provider", bestEntry.ProviderName,
		"tier", tier,
		"cost_per_min", bestEntry.CostPerMin,
		"off_peak", isOffPeak,
	)

	return ProviderCostProfile{
		ProviderName: bestEntry.ProviderName,
		CostPerMin:   bestEntry.CostPerMin,
		CostPerReq:   bestEntry.CostPerReq,
		Tier:         tier,
		FeatureFlags: bestEntry.FeatureFlags,
		FallbackChain: []string{bestEntry.ProviderName, "patter_secondary_fallback"},
	}
}

func (e *CostOptimizationEngine) SelectOptimalProvider(resType ResourceType, requiresPremiumFeatures bool) ProviderCostProfile {
	reqs := make(map[string]bool)
	if requiresPremiumFeatures {
		if resType == ResourceTypeMeeting {
			reqs["video_hd"] = true
		} else if resType == ResourceTypeVoice {
			reqs["krisp_noise_suppression"] = true
		} else if resType == ResourceTypeAvatar {
			reqs["realtime_webrtc"] = true
		}
	}
	return e.ResolveProvider(resType, reqs)
}
