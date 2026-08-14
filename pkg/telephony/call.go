package telephony

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"vocal-engine/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// OutboundCallRequest represents the JSON body accepted by the endpoint
type OutboundCallRequest struct {
	TargetNumber string `json:"targetNumber" binding:"required" example:"+15550199"`
	AgentPrompt  string `json:"agentPrompt" binding:"required" example:"You are a friendly assistant."`
	TenantID     string `json:"tenantId" binding:"required" example:"tenant-123"`
}

// OutboundCallResponse represents the success response
type OutboundCallResponse struct {
	CallSID string `json:"callSid" example:"CAXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`
	Status  string `json:"status" example:"queued"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request body"`
}

// CallController handles telephony-related HTTP endpoints
type CallController struct {
	cfg          *config.Config
	twilioClient *twilio.RestClient
}

// NewCallController creates a new CallController
func NewCallController(cfg *config.Config) *CallController {
	var restClient *twilio.RestClient
	if cfg.TwilioAccountSID != "" && cfg.TwilioAuthToken != "" {
		restClient = twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: cfg.TwilioAccountSID,
			Password: cfg.TwilioAuthToken,
		})
	} else {
		// Use default Twilio SDK client configuration (reads from environment variables)
		restClient = twilio.NewRestClient()
	}
	return &CallController{
		cfg:          cfg,
		twilioClient: restClient,
	}
}

// HandleOutboundCall triggers an outbound call with Twilio and connects it to the WebSocket media stream
// @Summary Trigger an outbound voice call
// @Description Initiates an outbound call via Twilio and connects the voice channel to the WebSocket stream
// @Tags Telephony
// @Accept json
// @Produce json
// @Param request body OutboundCallRequest true "Outbound Call Details"
// @Success 201 {object} OutboundCallResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/call/outbound [post]
func (cc *CallController) HandleOutboundCall(c *gin.Context) {
	var req OutboundCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("invalid request body: %v", err)})
		return;
	}

	// Dynamic callback URL which Twilio will invoke to fetch the TwiML.
	// Point to our TwiML generation handler.
	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	callbackURL := fmt.Sprintf("%s://%s/api/v1/call/twiml?prompt=%s&tenantId=%s",
		scheme, c.Request.Host, url.QueryEscape(req.AgentPrompt), url.QueryEscape(req.TenantID))

	params := &twilioApi.CreateCallParams{}
	params.SetTo(req.TargetNumber)
	params.SetFrom(cc.cfg.TwilioPhoneNumber)
	params.SetUrl(callbackURL)

	resp, err := cc.twilioClient.Api.CreateCall(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("failed to initiate Twilio call: %v", err)})
		return
	}

	var sid string
	var status string
	if resp.Sid != nil {
		sid = *resp.Sid
	}
	if resp.Status != nil {
		status = *resp.Status
	}

	c.JSON(http.StatusCreated, OutboundCallResponse{
		CallSID: sid,
		Status:  status,
	})
}

// TwiMLStream represents the TwiML XML structure for `<Connect><Stream>`
type TwiMLStream struct {
	XMLName xml.Name `xml:"Response"`
	Connect Connect  `xml:"Connect"`
}

type Connect struct {
	Stream Stream `xml:"Stream"`
}

type Stream struct {
	URL string `xml:"url,attr"`
}

// HandleTwiML generates the TwiML stream connection instruction for Twilio
func (cc *CallController) HandleTwiML(c *gin.Context) {
	prompt := c.Query("prompt")
	tenantID := c.Query("tenantId")

	// Websocket path for Twilio stream connection.
	// We pass down prompt and tenant ID through query parameters so the websocket stream handler has access to them.
	scheme := "wss"
	if c.Request.TLS == nil {
		scheme = "ws"
	}
	wsURL := fmt.Sprintf("%s://%s/ws/twilio/stream?prompt=%s&tenantId=%s",
		scheme, c.Request.Host, url.QueryEscape(prompt), url.QueryEscape(tenantID))

	twiml := TwiMLStream{
		Connect: Connect{
			Stream: Stream{
				URL: wsURL,
			},
		},
	}

	c.XML(http.StatusOK, twiml)
}
