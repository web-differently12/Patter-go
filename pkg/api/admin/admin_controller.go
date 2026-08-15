package admin

import (
	"net/http"

	"vocal-engine/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func NewAdminController() *AdminController {
	return &AdminController{}
}

// GetTenantConfig returns the active BYOK provider configuration for a tenant
// @Summary Get Tenant Configuration
// @Description Returns the active provider API keys for the tenant
// @Tags Admin - Tenant Configuration
// @Produce json
// @Param tenant_id path string true "Tenant Identifier"
// @Success 200 {object} dto.TenantConfigDTO
// @Failure 401 {object} telephony.ErrorResponse
// @Router /api/v1/admin/tenants/{tenant_id}/config [get]
func (a *AdminController) GetTenantConfig(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	config := dto.TenantConfigDTO{
		TenantID:             tenantID,
		TwilioAccountSID:     "AC_configured_for_" + tenantID,
		EvolutionGoServerURL: "http://evolution-go:8080",
		Webhooks: map[string]string{
			"voice": "https://gateway.example.com/api/v1/gateway/voice/inbound/webhook",
		},
	}
	c.JSON(http.StatusOK, config)
}

// UpdateTenantConfig updates BYOK provider API keys for a tenant
// @Summary Update Tenant Configuration
// @Description Updates active API keys and endpoints for a tenant
// @Tags Admin - Tenant Configuration
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant Identifier"
// @Param config body dto.TenantConfigDTO true "Tenant BYOK Configuration Specs"
// @Success 200 {object} dto.TenantConfigDTO
// @Failure 400 {object} telephony.ErrorResponse
// @Router /api/v1/admin/tenants/{tenant_id}/config [put]
func (a *AdminController) UpdateTenantConfig(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	var req dto.TenantConfigDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.TenantID = tenantID
	c.JSON(http.StatusOK, req)
}

// GetTenantBranding returns branding parameters for a tenant
// @Summary Get Tenant Branding
// @Description Returns custom domain, logo, display name, and voice default settings
// @Tags Admin - Tenant Configuration
// @Produce json
// @Param tenant_id path string true "Tenant Identifier"
// @Success 200 {object} dto.TenantBrandingDTO
// @Router /api/v1/admin/tenants/{tenant_id}/branding [get]
func (a *AdminController) GetTenantBranding(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	branding := dto.TenantBrandingDTO{
		TenantID:     tenantID,
		CustomDomain: "ai." + tenantID + ".com",
		DisplayName:  "AI Assistant",
		LogoURL:      "https://example.com/logo.png",
		DefaultVoice: "alloy",
		PrimaryColor: "#3B82F6",
	}
	c.JSON(http.StatusOK, branding)
}
