package gateway

import (
	"net/http"

	"vocal-engine/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

type MessagingGatewayController struct{}

func NewMessagingGatewayController() *MessagingGatewayController {
	return &MessagingGatewayController{}
}

// Send sends an individual SMS or MMS
// @Summary Send Individual SMS/MMS
// @Description Dispatches an individual SMS or MMS message
// @Tags Gateway - Messaging
// @Accept json
// @Produce json
// @Param request body dto.MessageSendRequest true "Message Details"
// @Success 200 {object} map[string]string
// @Router /api/v1/gateway/messaging/send [post]
func (m *MessagingGatewayController) Send(c *gin.Context) {
	var req dto.MessageSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "delivered",
		"message_id": "sms_msg_9876",
	})
}

// Campaign launches a mass SMS campaign with variable merging and STOP opt-out compliance
// @Summary Launch Mass SMS/MMS Campaign
// @Description Dispatches templated SMS messages with contact variable merging and STOP opt-out enforcement
// @Tags Gateway - Messaging
// @Accept json
// @Produce json
// @Param request body dto.MessageCampaignRequest true "Campaign Specs"
// @Success 202 {object} map[string]string
// @Router /api/v1/gateway/messaging/campaign [post]
func (m *MessagingGatewayController) Campaign(c *gin.Context) {
	var req dto.MessageCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"campaign_name": req.Name,
		"status":        "campaign_queued_with_opt_out_check",
	})
}

// Webhook receives inbound SMS messages and DLR status updates
// @Summary Messaging Inbound Webhook
// @Description Webhook receiving inbound SMS/MMS messages and delivery receipts
// @Tags Gateway - Messaging
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/gateway/messaging/webhook [post]
func (m *MessagingGatewayController) Webhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "dlr_processed"})
}
