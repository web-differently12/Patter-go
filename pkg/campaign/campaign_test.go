package campaign_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/campaign/controller"
	"github.com/lynxflow/patter-go/pkg/campaign/dto"
	"github.com/lynxflow/patter-go/pkg/campaign/service"
	"github.com/lynxflow/patter-go/pkg/core"
)

func TestAudienceResolver(t *testing.T) {
	svc := service.NewCampaignService()

	contacts := []dto.TargetContact{
		{
			ID:          "1",
			Name:        "Valid Contact AND",
			Phone:       "+15551234567",
			CategoryIDs: []string{"cat1"},
			Tags:        []string{"tag1", "vip"},
			OptedOut:    false,
		},
		{
			ID:          "2",
			Name:        "Opted Out Contact",
			Phone:       "+15559876543",
			CategoryIDs: []string{"cat1"},
			Tags:        []string{"tag1"},
			OptedOut:    true,
		},
		{
			ID:          "3",
			Name:        "Excluded Tag Contact",
			Phone:       "+15550001111",
			CategoryIDs: []string{"cat1"},
			Tags:        []string{"spam"},
			OptedOut:    false,
		},
		{
			ID:          "4",
			Name:        "Invalid Phone Contact",
			Phone:       "invalid_phone",
			CategoryIDs: []string{"cat1"},
			Tags:        []string{"tag1"},
			OptedOut:    false,
		},
	}

	filter := dto.AudienceFilter{
		MatchLogic:         "AND",
		ContactCategoryIDs: []string{"cat1"},
		RequiredTags:       []string{"tag1"},
		ExcludedTags:       []string{"spam"},
	}

	resolved := svc.ResolveAudience(filter, contacts)

	if len(resolved) != 1 {
		t.Fatalf("Expected 1 resolved contact, got %d", len(resolved))
	}
	if resolved[0].ID != "1" {
		t.Errorf("Expected contact ID 1, got %s", resolved[0].ID)
	}
}

func TestCalculateJitter(t *testing.T) {
	svc := service.NewCampaignService()

	baseDelay := 10
	for i := 0; i < 100; i++ {
		jitter := svc.CalculateJitter(baseDelay)
		if jitter < 0 {
			t.Fatalf("Jitter must be non-negative, got %v", jitter)
		}
		// Jitter should be within +/- 30% of baseDelay (i.e. between 7s and 13s)
		if jitter < 7*time.Second || jitter > 13*time.Second {
			t.Errorf("Jitter out of expected bounds [7s, 13s]: got %v", jitter)
		}
	}
}

func TestInterpolate(t *testing.T) {
	svc := service.NewCampaignService()

	template := "Hello {{first_name}}, welcome to {{company}}!"
	vars := map[string]string{
		"first_name": "Jean",
		"company":    "Lynxflow",
	}

	result := svc.Interpolate(template, vars)
	expected := "Hello Jean, welcome to Lynxflow!"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestCampaignLifecycleREST(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := service.NewCampaignService()
	ctrl := controller.NewCampaignController(svc)

	r := gin.New()
	r.Use(core.AuthMiddleware())
	r.POST("/api/v1/gateway/campaigns", ctrl.CreateCampaign)
	r.GET("/api/v1/gateway/campaigns", ctrl.ListCampaigns)
	r.GET("/api/v1/gateway/campaigns/:id", ctrl.GetCampaign)
	r.POST("/api/v1/gateway/campaigns/:id/pause", ctrl.PauseCampaign)
	r.POST("/api/v1/gateway/campaigns/:id/resume", ctrl.ResumeCampaign)

	// 1. Create Campaign
	body := `{
		"name": "Test Campaign",
		"channel": "WHATSAPP",
		"session_names": ["sess1", "sess2"],
		"message_content": "Test content {{first_name}}",
		"random_delay_sec": 0,
		"start_immediately": true
	}`

	req, _ := http.NewRequest("POST", "/api/v1/gateway/campaigns", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant_test")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// 2. List Campaigns
	reqList, _ := http.NewRequest("GET", "/api/v1/gateway/campaigns", nil)
	reqList.Header.Set("X-Tenant-ID", "tenant_test")

	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for list, got %d", wList.Code)
	}
}

func TestVoiceCampaignREST(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := service.NewCampaignService()
	ctrl := controller.NewCampaignController(svc)

	r := gin.New()
	r.Use(core.AuthMiddleware())
	r.POST("/api/v1/gateway/voice/campaign", ctrl.CreateVoiceCampaign)

	body := `{
		"campaign_name": "Test Voice Campaign",
		"from_numbers": ["+15550001111"],
		"prompt_template": "Bonjour {{phone}}, vous avez un message vocal.",
		"random_delay_sec": 0,
		"target_phone_numbers": ["+15552223333", "+15554445555"]
	}`

	req, _ := http.NewRequest("POST", "/api/v1/gateway/voice/campaign", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant_test")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for voice campaign creation, got %d", w.Code)
	}
}
