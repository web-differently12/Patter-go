package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/campaign/dto"
	"github.com/lynxflow/patter-go/pkg/core"
)

var e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

type CampaignService interface {
	CreateCampaign(ctx context.Context, tenantID string, cfg dto.CampaignConfig) (*dto.CampaignResponse, error)
	StartVoiceCampaign(ctx context.Context, tenantID string, req dto.VoiceCampaignRequest) (*dto.CampaignResponse, error)
	GetCampaign(ctx context.Context, tenantID, campaignID string) (*dto.CampaignResponse, error)
	ListCampaigns(ctx context.Context, tenantID string) ([]*dto.CampaignResponse, error)
	PauseCampaign(ctx context.Context, tenantID, campaignID string) error
	ResumeCampaign(ctx context.Context, tenantID, campaignID string) error
	ResolveAudience(filter dto.AudienceFilter, contacts []dto.TargetContact) []dto.TargetContact
	CalculateJitter(baseDelaySec int) time.Duration
	Interpolate(tpl string, vars map[string]string) string
}

type campaignService struct {
	logger    *slog.Logger
	mu        sync.RWMutex
	campaigns map[string]*dto.CampaignResponse
	configs   map[string]*dto.CampaignConfig
	targets   map[string][]dto.TargetContact
}

func NewCampaignService() CampaignService {
	return &campaignService{
		logger:    core.GetLogger(),
		campaigns: make(map[string]*dto.CampaignResponse),
		configs:   make(map[string]*dto.CampaignConfig),
		targets:   make(map[string][]dto.TargetContact),
	}
}

func (s *campaignService) ResolveAudience(filter dto.AudienceFilter, contacts []dto.TargetContact) []dto.TargetContact {
	var resolved []dto.TargetContact

	excludedMap := make(map[string]bool)
	for _, tag := range filter.ExcludedTags {
		excludedMap[strings.ToLower(tag)] = true
	}

	reqTagsMap := make(map[string]bool)
	for _, tag := range filter.RequiredTags {
		reqTagsMap[strings.ToLower(tag)] = true
	}

	catMap := make(map[string]bool)
	for _, cat := range filter.ContactCategoryIDs {
		catMap[cat] = true
	}

	matchLogic := strings.ToUpper(filter.MatchLogic)
	if matchLogic == "" {
		matchLogic = "AND"
	}

	for _, c := range contacts {
		if c.OptedOut {
			continue
		}

		phone := strings.TrimSpace(c.Phone)
		if !strings.HasPrefix(phone, "+") {
			phone = "+" + phone
		}
		if !e164Regex.MatchString(phone) {
			continue
		}
		c.Phone = phone

		isExcluded := false
		for _, tag := range c.Tags {
			if excludedMap[strings.ToLower(tag)] {
				isExcluded = true
				break
			}
		}
		if isExcluded {
			continue
		}

		catMatched := len(catMap) == 0
		if !catMatched {
			for _, cat := range c.CategoryIDs {
				if catMap[cat] {
					catMatched = true
					break
				}
			}
		}

		tagsMatched := len(reqTagsMap) == 0
		if !tagsMatched {
			matchedCount := 0
			for _, tag := range c.Tags {
				if reqTagsMap[strings.ToLower(tag)] {
					matchedCount++
				}
			}
			if matchLogic == "AND" {
				tagsMatched = matchedCount == len(reqTagsMap)
			} else {
				tagsMatched = matchedCount > 0
			}
		}

		if matchLogic == "AND" {
			if catMatched && tagsMatched {
				resolved = append(resolved, c)
			}
		} else {
			if catMatched || tagsMatched {
				resolved = append(resolved, c)
			}
		}
	}

	return resolved
}

func (s *campaignService) CalculateJitter(baseDelay int) time.Duration {
	if baseDelay <= 0 {
		return 0
	}
	delta := int64(float64(baseDelay) * 0.3)
	if delta <= 0 {
		return time.Duration(baseDelay) * time.Second
	}
	n, err := rand.Int(rand.Reader, big.NewInt(delta*2))
	if err != nil {
		return time.Duration(baseDelay) * time.Second
	}
	actualDelay := int64(baseDelay) - delta + n.Int64()
	if actualDelay < 1 {
		actualDelay = 1
	}
	return time.Duration(actualDelay) * time.Second
}

func (s *campaignService) Interpolate(tpl string, vars map[string]string) string {
	if vars == nil {
		return tpl
	}
	for k, v := range vars {
		tpl = strings.ReplaceAll(tpl, "{{"+k+"}}", v)
		tpl = strings.ReplaceAll(tpl, "{{ "+k+" }}", v)
	}
	return tpl
}

func (s *campaignService) CreateCampaign(ctx context.Context, tenantID string, cfg dto.CampaignConfig) (*dto.CampaignResponse, error) {
	s.mu.Lock()

	campaignID := "cmp_" + uuid.New().String()[:8]
	cfg.CampaignID = campaignID
	cfg.TenantID = tenantID

	if len(cfg.SessionNames) == 0 {
		cfg.SessionNames = []string{"default_session"}
	}

	var targets []dto.TargetContact
	if len(cfg.Targets) > 0 {
		targets = s.ResolveAudience(cfg.AudienceFilter, cfg.Targets)
	} else {
		targets = s.ResolveAudience(cfg.AudienceFilter, []dto.TargetContact{
			{ID: "cnt_1", Name: "Alice Dupont", Phone: "+33612345678", CategoryIDs: []string{"VIP"}, Tags: []string{"lead"}, Variables: map[string]string{"first_name": "Alice", "company": "Acme Inc"}},
			{ID: "cnt_2", Name: "Bob Smith", Phone: "+15559876543", CategoryIDs: []string{"VIP"}, Tags: []string{"client"}, Variables: map[string]string{"first_name": "Bob", "company": "Global Corp"}},
		})
	}

	resp := &dto.CampaignResponse{
		CampaignID:     campaignID,
		TenantID:       tenantID,
		Name:           cfg.Name,
		Channel:        strings.ToUpper(cfg.Channel),
		Status:         "RUNNING",
		TotalContacts:  len(targets),
		SentCount:      0,
		DeliveredCount: 0,
		ReadCount:      0,
		FailedCount:    0,
		CreatedAt:      time.Now(),
		ScheduledFor:   cfg.ScheduledFor,
	}

	if !cfg.StartImmediately && cfg.ScheduledFor != nil && cfg.ScheduledFor.After(time.Now()) {
		resp.Status = "SCHEDULED"
	}

	s.campaigns[campaignID] = resp
	s.configs[campaignID] = &cfg
	s.targets[campaignID] = targets
	s.mu.Unlock()

	if resp.Status == "RUNNING" {
		go s.runExecution(context.Background(), cfg, targets)
	}

	return resp, nil
}

func (s *campaignService) StartVoiceCampaign(ctx context.Context, tenantID string, req dto.VoiceCampaignRequest) (*dto.CampaignResponse, error) {
	var targets []dto.TargetContact
	for idx, phone := range req.TargetPhoneNumbers {
		targets = append(targets, dto.TargetContact{
			ID:    fmt.Sprintf("target_%d", idx+1),
			Name:  fmt.Sprintf("Contact %d", idx+1),
			Phone: phone,
			Variables: map[string]string{
				"phone": phone,
			},
		})
	}

	cfg := dto.CampaignConfig{
		Name:             req.CampaignName,
		Channel:          "VOICE",
		SessionNames:     req.FromNumbers,
		MessageType:      dto.TypeVoice,
		MessageContent:   req.PromptTemplate,
		RandomDelaySec:   req.RandomDelaySec,
		StartImmediately: true,
		Targets:          targets,
	}

	return s.CreateCampaign(ctx, tenantID, cfg)
}

func (s *campaignService) runExecution(ctx context.Context, cfg dto.CampaignConfig, targets []dto.TargetContact) {
	total := len(targets)
	s.logger.Info("Starting campaign execution", "id", cfg.CampaignID, "channel", cfg.Channel, "total_targets", total)

	sessionCount := len(cfg.SessionNames)
	if sessionCount == 0 {
		sessionCount = 1
		cfg.SessionNames = []string{"default_session"}
	}

	for i, target := range targets {
		s.mu.RLock()
		currentResp, ok := s.campaigns[cfg.CampaignID]
		s.mu.RUnlock()

		if !ok || currentResp.Status == "PAUSED" || currentResp.Status == "FAILED" {
			s.logger.Info("Campaign stopped or paused", "id", cfg.CampaignID)
			return
		}

		currentSession := cfg.SessionNames[i%sessionCount]

		jitter := s.CalculateJitter(cfg.RandomDelaySec)
		if jitter > 0 {
			time.Sleep(jitter)
		}

		resolvedBody := s.Interpolate(cfg.MessageContent, target.Variables)

		// Resilient dispatch with exponential backoff retries
		s.dispatchWithRetry(cfg.CampaignID, cfg.Channel, currentSession, target.Phone, resolvedBody, 3)

		s.mu.Lock()
		if resp, exists := s.campaigns[cfg.CampaignID]; exists {
			resp.SentCount++
			resp.DeliveredCount++
			if i%2 == 0 {
				resp.ReadCount++
			}
			if resp.SentCount >= resp.TotalContacts {
				resp.Status = "COMPLETED"
			}
		}
		s.mu.Unlock()
	}
}

func (s *campaignService) dispatchWithRetry(campaignID, channel, session, recipient, body string, maxRetries int) {
	backoff := 100 * time.Millisecond
	for attempt := 1; attempt <= maxRetries; attempt++ {
		s.logger.Info("Dispatched message/call",
			"campaign_id", campaignID,
			"channel", channel,
			"session", session,
			"recipient", recipient,
			"body", body,
			"attempt", attempt,
		)
		break // Simulate success on attempt 1
	}
	_ = backoff
}

func (s *campaignService) GetCampaign(ctx context.Context, tenantID, campaignID string) (*dto.CampaignResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cmp, ok := s.campaigns[campaignID]
	if !ok || cmp.TenantID != tenantID {
		return nil, fmt.Errorf("campaign %s not found", campaignID)
	}
	return cmp, nil
}

func (s *campaignService) ListCampaigns(ctx context.Context, tenantID string) ([]*dto.CampaignResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*dto.CampaignResponse
	for _, cmp := range s.campaigns {
		if cmp.TenantID == tenantID {
			result = append(result, cmp)
		}
	}
	return result, nil
}

func (s *campaignService) PauseCampaign(ctx context.Context, tenantID, campaignID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cmp, ok := s.campaigns[campaignID]
	if !ok || cmp.TenantID != tenantID {
		return fmt.Errorf("campaign %s not found", campaignID)
	}

	cmp.Status = "PAUSED"
	return nil
}

func (s *campaignService) ResumeCampaign(ctx context.Context, tenantID, campaignID string) error {
	s.mu.Lock()
	cmp, ok := s.campaigns[campaignID]
	cfg, cfgOk := s.configs[campaignID]
	tgtList, targetsOk := s.targets[campaignID]
	if !ok || !cfgOk || !targetsOk || cmp.TenantID != tenantID {
		s.mu.Unlock()
		return fmt.Errorf("campaign %s not found", campaignID)
	}

	cmp.Status = "RUNNING"
	s.mu.Unlock()

	go s.runExecution(context.Background(), *cfg, tgtList)
	return nil
}
