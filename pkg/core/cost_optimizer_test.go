package core_test

import (
	"testing"

	"github.com/lynxflow/patter-go/pkg/core"
)

func TestCostOptimizationEngine(t *testing.T) {
	engine := core.NewCostOptimizationEngine(nil)

	// Test Meeting Cost Routing
	economyMeeting := engine.SelectOptimalProvider(core.ResourceTypeMeeting, false)
	if economyMeeting.ProviderName != "patter_native_rtc" {
		t.Errorf("Expected patter_native_rtc for non-premium meeting, got %s", economyMeeting.ProviderName)
	}

	premiumMeeting := engine.SelectOptimalProvider(core.ResourceTypeMeeting, true)
	if premiumMeeting.ProviderName != "patter_premium_engine" {
		t.Errorf("Expected patter_premium_engine for premium meeting, got %s", premiumMeeting.ProviderName)
	}

	// Test Voice Cost Routing
	economyVoice := engine.SelectOptimalProvider(core.ResourceTypeVoice, false)
	if economyVoice.ProviderName != "patter_sip_direct_pbx" {
		t.Errorf("Expected patter_sip_direct_pbx for economy voice, got %s", economyVoice.ProviderName)
	}
}
