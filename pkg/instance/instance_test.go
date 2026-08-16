package instance_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
	"github.com/lynxflow/patter-go/pkg/instance/controller"
	"github.com/lynxflow/patter-go/pkg/instance/service"
)

func TestInstanceLifecycleREST(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := service.NewInstanceService()
	ctrl := controller.NewInstanceController(svc)

	r := gin.New()
	r.Use(core.AuthMiddleware())
	r.POST("/api/v1/gateway/instances", ctrl.CreateInstance)
	r.GET("/api/v1/gateway/instances", ctrl.ListInstances)

	// Create Instance
	body := `{"instance_name": "Tenant Acme Instance", "description": "Instance White Label"}`
	req, _ := http.NewRequest("POST", "/api/v1/gateway/instances", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant_acme")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	// List Instances
	reqList, _ := http.NewRequest("GET", "/api/v1/gateway/instances", nil)
	reqList.Header.Set("X-Tenant-ID", "tenant_acme")

	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", wList.Code)
	}
}
