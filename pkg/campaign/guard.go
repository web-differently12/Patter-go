package campaign

import (
	"log/slog"
	"strings"

	"github.com/lynxflow/patter-go/pkg/campaign/dto"
	"github.com/lynxflow/patter-go/pkg/core"
)

type AdminCompliancePolicy struct {
	MaxDispatchesPerHour int      `json:"max_dispatches_per_hour"`
	MinJitterSec         int      `json:"min_jitter_sec"`
	RequireOptOutLink    bool     `json:"require_opt_out_link"`
	ForbiddenKeywords    []string `json:"forbidden_keywords"`
	StrictE164Validation bool     `json:"strict_e164_validation"`
}

type ComplianceGuardResult struct {
	Compliant       bool               `json:"compliant"`
	AdjustedConfig  dto.CampaignConfig `json:"adjusted_config"`
	Violations      []string           `json:"violations"`
	AutoCorrected   bool               `json:"auto_corrected"`
}

type CampaignComplianceGuard struct {
	policy AdminCompliancePolicy
	logger *slog.Logger
}

func NewCampaignComplianceGuard(policy AdminCompliancePolicy, logger *slog.Logger) *CampaignComplianceGuard {
	if logger == nil {
		logger = core.GetLogger()
	}
	if policy.MinJitterSec <= 0 {
		policy.MinJitterSec = 10
	}
	if policy.MaxDispatchesPerHour <= 0 {
		policy.MaxDispatchesPerHour = 500
	}
	if len(policy.ForbiddenKeywords) == 0 {
		policy.ForbiddenKeywords = []string{"gagnez 1000000€", "urgent virement", "casino sans depot"}
	}
	return &CampaignComplianceGuard{
		policy: policy,
		logger: logger.With("component", "campaign_compliance_guard"),
	}
}

func (g *CampaignComplianceGuard) InspectAndAdjust(cfg dto.CampaignConfig) ComplianceGuardResult {
	var violations []string
	adjusted := cfg
	autoCorrected := false

	// 1. Minimum Anti-Spam Jitter Check
	if cfg.RandomDelaySec < g.policy.MinJitterSec {
		violations = append(violations, "Jitter delay too low, risked operator blacklisting")
		adjusted.RandomDelaySec = g.policy.MinJitterSec
		autoCorrected = true
	}

	// 2. Forbidden Keywords Filter
	for _, kw := range g.policy.ForbiddenKeywords {
		if strings.Contains(strings.ToLower(cfg.MessageContent), strings.ToLower(kw)) {
			violations = append(violations, "Forbidden keyword detected: "+kw)
			adjusted.MessageContent = strings.ReplaceAll(adjusted.MessageContent, kw, "[contenu modéré]")
			autoCorrected = true
		}
	}

	// 3. Opt-Out Mention Enforcement for Marketing
	if g.policy.RequireOptOutLink && cfg.Channel == "WHATSAPP" {
		if !strings.Contains(strings.ToLower(cfg.MessageContent), "stop") {
			violations = append(violations, "Missing mandatory opt-out mention")
			adjusted.MessageContent += " (Répondez STOP pour vous désinscrire)"
			autoCorrected = true
		}
	}

	if len(violations) > 0 {
		g.logger.Warn("Campaign compliance violations detected & auto-corrected to prevent blacklisting",
			"campaign_id", cfg.CampaignID,
			"tenant_id", cfg.TenantID,
			"violations_count", len(violations),
		)
	}

	return ComplianceGuardResult{
		Compliant:      len(violations) == 0,
		AdjustedConfig: adjusted,
		Violations:     violations,
		AutoCorrected:  autoCorrected,
	}
}
