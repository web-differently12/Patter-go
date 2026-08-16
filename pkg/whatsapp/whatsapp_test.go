package whatsapp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
	"github.com/lynxflow/patter-go/pkg/whatsapp/controller"
	"github.com/lynxflow/patter-go/pkg/whatsapp/service"
)

func TestWhatsAppProxyREST(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := service.NewWhatsAppService()
	ctrl := controller.NewWhatsAppController(svc)

	r := gin.New()
	r.Use(core.AuthMiddleware())
	r.POST("/api/v1/gateway/whatsapp/connect", ctrl.ConnectSession)
	r.GET("/api/v1/gateway/whatsapp/qrcode", ctrl.GetQRCode)
	r.GET("/api/v1/gateway/whatsapp/sessions", ctrl.ListSessions)
	r.POST("/api/v1/gateway/whatsapp/message/send", ctrl.SendMessage)

	// 1. Connect session
	body := `{"session_name": "test_sess_1", "phone_number": "+15550001234"}`
	req, _ := http.NewRequest("POST", "/api/v1/gateway/whatsapp/connect", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant_wa_test")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Get QR Code
	reqQR, _ := http.NewRequest("GET", "/api/v1/gateway/whatsapp/qrcode?session=test_sess_1", nil)
	reqQR.Header.Set("X-Tenant-ID", "tenant_wa_test")

	wQR := httptest.NewRecorder()
	r.ServeHTTP(wQR, reqQR)

	if wQR.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for QR Code, got %d", wQR.Code)
	}

	// 3. Send Message
	sendBody := `{"session_name": "test_sess_1", "recipient": "+15559998888", "message": "Test message"}`
	reqSend, _ := http.NewRequest("POST", "/api/v1/gateway/whatsapp/message/send", strings.NewReader(sendBody))
	reqSend.Header.Set("Content-Type", "application/json")
	reqSend.Header.Set("X-Tenant-ID", "tenant_wa_test")

	wSend := httptest.NewRecorder()
	r.ServeHTTP(wSend, reqSend)

	if wSend.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for Send Message, got %d", wSend.Code)
	}
}
