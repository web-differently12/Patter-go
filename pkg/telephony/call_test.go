package telephony

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"vocal-engine/pkg/config"

	"github.com/gin-gonic/gin"
)

func TestHandleTwiML(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		Port: "8080",
	}

	controller := NewCallController(cfg)
	r := gin.New()
	r.GET("/api/v1/call/twiml", controller.HandleTwiML)

	req, err := http.NewRequest("GET", "/api/v1/call/twiml?prompt=test-prompt&tenantId=test-tenant-123", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", w.Code)
	}

	// Unmarshal XML response
	var response TwiMLStream
	if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal XML response: %v. Body was: %s", err, w.Body.String())
	}

	expectedURL := "ws:///ws/twilio/stream?prompt=test-prompt&tenantId=test-tenant-123"
	if response.Connect.Stream.URL != expectedURL {
		t.Errorf("expected Stream URL %q, got %q", expectedURL, response.Connect.Stream.URL)
	}
}
