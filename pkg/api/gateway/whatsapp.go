package gateway

import (
	"net/http"

	"vocal-engine/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

type WhatsAppGatewayController struct{}

func NewWhatsAppGatewayController() *WhatsAppGatewayController {
	return &WhatsAppGatewayController{}
}

// Send sends an individual WhatsApp message
// @Summary Send Individual WhatsApp Message
// @Description Routes individual WhatsApp message for the client tenant
// @Tags Gateway - WhatsApp
// @Accept json
// @Produce json
// @Param request body dto.WhatsAppSendRequest true "Message Details"
// @Success 200 {object} map[string]string
// @Router /api/v1/gateway/whatsapp/send [post]
func (w *WhatsAppGatewayController) Send(c *gin.Context) {
	var req dto.WhatsAppSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "sent",
		"message_id": "wa_msg_12345",
	})
}

// Campaign launches a mass WhatsApp campaign using approved templates
// @Summary Launch Mass WhatsApp Campaign
// @Description Dispatches template-approved WhatsApp messages to a list of recipients
// @Tags Gateway - WhatsApp
// @Accept json
// @Produce json
// @Param request body dto.WhatsAppCampaignRequest true "Campaign Specs"
// @Success 202 {object} map[string]string
// @Router /api/v1/gateway/whatsapp/campaign [post]
func (w *WhatsAppGatewayController) Campaign(c *gin.Context) {
	var req dto.WhatsAppCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"campaign_name": req.CampaignName,
		"template_used": req.TemplateName,
		"status":        "processing_hsm_campaign",
	})
}

// GetTemplates lists approved WhatsApp templates
// @Summary List Approved WhatsApp Templates
// @Description Returns approved message templates for the tenant's account
// @Tags Gateway - WhatsApp
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/gateway/whatsapp/templates [get]
func (w *WhatsAppGatewayController) GetTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "approved_templates_retrieved"})
}

// Webhook receives incoming WhatsApp messages
// @Summary WhatsApp Inbound Webhook
// @Description Webhook endpoint receiving inbound WhatsApp messages and status updates
// @Tags Gateway - WhatsApp
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/gateway/whatsapp/webhook [post]
func (w *WhatsAppGatewayController) Webhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "event_received"})
}
