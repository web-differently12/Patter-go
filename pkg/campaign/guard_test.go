package campaign_test

import (
	"testing"

	"github.com/lynxflow/patter-go/pkg/campaign"
	"github.com/lynxflow/patter-go/pkg/campaign/dto"
)

func TestCampaignComplianceGuard(t *testing.T) {
	policy := campaign.AdminCompliancePolicy{
		MinJitterSec:      15,
		RequireOptOutLink: true,
		ForbiddenKeywords: []string{"gagnez 1000000€"},
	}

	guard := campaign.NewCampaignComplianceGuard(policy, nil)

	unsafeCfg := dto.CampaignConfig{
		CampaignID:     "cmp_unsafe_01",
		Channel:        "WHATSAPP",
		RandomDelaySec: 2, // Violates MinJitterSec (15s)
		MessageContent: "Offre exceptionnelle ! gagnez 1000000€ tout de suite !",
	}

	res := guard.InspectAndAdjust(unsafeCfg)

	if res.Compliant {
		t.Errorf("Expected Compliant == false for unsafe config")
	}
	if !res.AutoCorrected {
		t.Errorf("Expected AutoCorrected == true")
	}
	if res.AdjustedConfig.RandomDelaySec != 15 {
		t.Errorf("Expected AdjustedConfig.RandomDelaySec == 15, got %d", res.AdjustedConfig.RandomDelaySec)
	}
	if res.AdjustedConfig.MessageContent == unsafeCfg.MessageContent {
		t.Errorf("Expected forbidden keyword to be replaced in AdjustedConfig")
	}
}
