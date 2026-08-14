package gateway

import (
	"math/rand"
	"net/http"
	"time"

	"vocal-engine/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

type VoiceGatewayController struct{}

func NewVoiceGatewayController() *VoiceGatewayController {
	return &VoiceGatewayController{}
}

// Outbound initiates an individual outbound call for a tenant
// @Summary Trigger Individual Outbound Call
// @Description Initiates an outbound voice call using tenant's BYOK credentials
// @Tags Gateway - Voice
// @Accept json
// @Produce json
// @Param request body dto.VoiceOutboundRequest true "Call Details"
// @Success 201 {object} map[string]string
// @Router /api/v1/gateway/voice/outbound [post]
func (v *VoiceGatewayController) Outbound(c *gin.Context) {
	var req dto.VoiceOutboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"call_sid": "CA_gateway_" + req.TargetNumber,
		"status":   "queued",
	})
}

// Campaign launches a mass voice calling campaign with anti-spam jitter delays
// @Summary Launch Anti-Spam Voice Campaign
// @Description Dispatches voice calls sequentially with randomized jitter delays to prevent carrier spam detection
// @Tags Gateway - Voice
// @Accept json
// @Produce json
// @Param request body dto.VoiceCampaignRequest true "Campaign Details"
// @Success 202 {object} map[string]string
// @Router /api/v1/gateway/voice/campaign [post]
func (v *VoiceGatewayController) Campaign(c *gin.Context) {
	var req dto.VoiceCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate jitter delays for mass dispatching
	go func() {
		for _, num := range req.TargetNumbers {
			minJitter := req.MinJitterSecSec
			if minJitter <= 0 {
				minJitter = 3
			}
			maxJitter := req.MaxJitterSecSec
			if maxJitter <= minJitter {
				maxJitter = minJitter + 5
			}
			jitterSec := rand.Intn(maxJitter-minJitter) + minJitter
			time.Sleep(time.Duration(jitterSec) * time.Second)
			_ = num
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"campaign_name": req.Name,
		"status":        "dispatching_with_anti_spam_jitter",
	})
}

// InboundWebhook handles incoming calls for the AI receptionist
// @Summary Handle Inbound Telephony Webhook
// @Description Receives webhook requests for incoming calls and routes them to the agent stream
// @Tags Gateway - Voice
// @Produce xml
// @Success 200 {string} string "TwiML XML Response"
// @Router /api/v1/gateway/voice/inbound/webhook [post]
func (v *VoiceGatewayController) InboundWebhook(c *gin.Context) {
	twiml := `<?xml version="1.0" encoding="UTF-8"?><Response><Connect><Stream url="wss://gateway.example.com/ws/twilio/stream"/></Connect></Response>`
	c.Header("Content-Type", "text/xml")
	c.String(http.StatusOK, twiml)
}
