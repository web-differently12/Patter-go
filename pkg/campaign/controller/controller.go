package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/campaign/dto"
	"github.com/lynxflow/patter-go/pkg/campaign/service"
	"github.com/lynxflow/patter-go/pkg/core"
)

type CampaignController struct {
	campaignService service.CampaignService
}

func NewCampaignController(svc service.CampaignService) *CampaignController {
	return &CampaignController{campaignService: svc}
}

// CreateCampaign godoc
// @Summary      Lancer une campagne omnicanale (WhatsApp, SMS, Voice)
// @Description  Déclenche une campagne avec cadenceur anti-spam jitter et rotation de sessions
// @Tags         Gateway - Campaigns
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.CampaignConfig true "Configuration de la campagne"
// @Success      200 {object} core.APIResponse{data=dto.CampaignResponse}
// @Failure      400 {object} dto.ErrorResponse
// @Router       /api/v1/gateway/campaigns [post]
func (ctrl *CampaignController) CreateCampaign(c *gin.Context) {
	var req dto.CampaignConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.campaignService.CreateCampaign(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	core.Success(c, resp)
}

// CreateVoiceCampaign godoc
// @Summary      Lancer une campagne vocale avec Jitter Anti-Spam
// @Description  Déclenche une campagne d'appels de masse avec cadenceur aléatoire
// @Tags         Gateway - Voice Campaigns
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.VoiceCampaignRequest true "Configuration de la campagne"
// @Success      200 {object} core.APIResponse{data=dto.CampaignResponse}
// @Failure      400 {object} dto.ErrorResponse
// @Router       /api/v1/gateway/voice/campaign [post]
func (ctrl *CampaignController) CreateVoiceCampaign(c *gin.Context) {
	var req dto.VoiceCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.campaignService.StartVoiceCampaign(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	core.Success(c, resp)
}

// ListCampaigns godoc
// @Summary      Lister les campagnes du Tenant
// @Tags         Gateway - Campaigns
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]dto.CampaignResponse}
// @Router       /api/v1/gateway/campaigns [get]
func (ctrl *CampaignController) ListCampaigns(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.campaignService.ListCampaigns(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	core.Success(c, resp)
}

// GetCampaign godoc
// @Summary      Obtenir les détails d'une campagne
// @Tags         Gateway - Campaigns
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        id path string true "ID de la campagne"
// @Success      200 {object} core.APIResponse{data=dto.CampaignResponse}
// @Router       /api/v1/gateway/campaigns/{id} [get]
func (ctrl *CampaignController) GetCampaign(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	campaignID := c.Param("id")

	resp, err := ctrl.campaignService.GetCampaign(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	core.Success(c, resp)
}

// PauseCampaign godoc
// @Summary      Mettre en pause une campagne
// @Tags         Gateway - Campaigns
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        id path string true "ID de la campagne"
// @Success      200 {object} core.APIResponse
// @Router       /api/v1/gateway/campaigns/{id}/pause [post]
func (ctrl *CampaignController) PauseCampaign(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	campaignID := c.Param("id")

	if err := ctrl.campaignService.PauseCampaign(c.Request.Context(), tenantID, campaignID); err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	core.Success(c, gin.H{"message": "Campaign paused successfully"})
}

// ResumeCampaign godoc
// @Summary      Reprendre une campagne en pause
// @Tags         Gateway - Campaigns
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        id path string true "ID de la campagne"
// @Success      200 {object} core.APIResponse
// @Router       /api/v1/gateway/campaigns/{id}/resume [post]
func (ctrl *CampaignController) ResumeCampaign(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	campaignID := c.Param("id")

	if err := ctrl.campaignService.ResumeCampaign(c.Request.Context(), tenantID, campaignID); err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	core.Success(c, gin.H{"message": "Campaign resumed successfully"})
}
