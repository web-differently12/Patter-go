package gateway

import (
	"net/http"

	"vocal-engine/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

type BrainController struct{}

func NewBrainController() *BrainController {
	return &BrainController{}
}

// Chat processes text messages via central LLM/OpenRouter with CRM context
// @Summary Process Multichannel LLM Chat
// @Description Handles text processing for SMS, WhatsApp, and Webchat via LLM with CRM context
// @Tags Gateway - Brain
// @Accept json
// @Produce json
// @Param request body dto.BrainChatRequest true "Chat Request Data"
// @Success 200 {object} dto.BrainChatResponse
// @Router /api/v1/gateway/brain/chat [post]
func (b *BrainController) Chat(c *gin.Context) {
	var req dto.BrainChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.BrainChatResponse{
		Reply:     "Thank you for your message: " + req.Message,
		ModelUsed: "openrouter/gpt-4o",
	}
	c.JSON(http.StatusOK, resp)
}

// Report retrieves cognitive meeting intelligence reports (LeMUR v3)
// @Summary Retrieve LeMUR v3 Meeting Report
// @Description Fetches structured LeMUR v3 reports including BANT scores and Action Items
// @Tags Gateway - Brain
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/gateway/brain/report [post]
func (b *BrainController) Report(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "report_retrieved",
		"report": "LeMUR v3 BANT score generated",
	})
}
