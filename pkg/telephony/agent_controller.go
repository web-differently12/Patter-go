package telephony

import (
	"net/http"
	"net/url"
	"fmt"

	"vocal-engine/pkg/engine"

	"github.com/gin-gonic/gin"
)

// AgentController manages Agent CRUD operations and inbound telephony routing
type AgentController struct {
	repo engine.AgentRepository
}

// NewAgentController creates a new AgentController
func NewAgentController(repo engine.AgentRepository) *AgentController {
	return &AgentController{repo: repo}
}

// HandleCreateAgent creates a new configured agent profile
// @Summary Create an AI Agent profile
// @Description Configures an AI agent's voice, prompts, tools, STT/TTS engine, and failover fallbacks
// @Tags Agents
// @Accept json
// @Produce json
// @Param request body engine.Agent true "Agent Profile Specifications"
// @Success 201 {object} engine.Agent
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/agents [post]
func (ac *AgentController) HandleCreateAgent(c *gin.Context) {
	var agent engine.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if agent.ID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "agent ID is required"})
		return
	}

	err := ac.repo.CreateAgent(c, &agent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// HandleListAgents lists all configured voice agents
// @Summary List configured AI Agents
// @Description Retrieves all voice agent profiles currently registered in the system
// @Tags Agents
// @Produce json
// @Success 200 {array} engine.Agent
// @Router /api/v1/agents [get]
func (ac *AgentController) HandleListAgents(c *gin.Context) {
	agents, err := ac.repo.ListAgents(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, agents)
}

// HandleInboundCall handles Twilio/Carrier incoming webhooks to connect calls to agents
// @Summary Handle Inbound Telephony Webhooks
// @Description Receives webhook requests from Twilio for incoming calls and connects them to a configured agent stream
// @Tags Telephony
// @Produce xml
// @Param agentId query string true "Registered Agent ID"
// @Param tenantId query string true "Registered Tenant ID"
// @Success 200 {string} string "TwiML XML response"
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/call/inbound [post]
func (ac *AgentController) HandleInboundCall(c *gin.Context) {
	agentID := c.Query("agentId")
	tenantID := c.Query("tenantId")

	agent, err := ac.repo.GetAgentByID(c, agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("configured agent %s not found: %v", agentID, err)})
		return
	}

	// Connect inbound call to custom agent stream WebSocket endpoint
	scheme := "wss"
	if c.Request.TLS == nil {
		scheme = "ws"
	}
	wsURL := fmt.Sprintf("%s://%s/ws/twilio/stream?prompt=%s&tenantId=%s&agentId=%s",
		scheme, c.Request.Host, url.QueryEscape(agent.SystemPrompt), url.QueryEscape(tenantID), url.QueryEscape(agentID))

	// Generate stream xml
	twiml := TwiMLStream{
		Connect: Connect{
			Stream: Stream{
				URL: wsURL,
			},
		},
	}

	c.XML(http.StatusOK, twiml)
}
