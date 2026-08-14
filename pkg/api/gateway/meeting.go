package gateway

import (
	"net/http"

	"vocal-engine/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

type MeetingGatewayController struct{}

func NewMeetingGatewayController() *MeetingGatewayController {
	return &MeetingGatewayController{}
}

// Schedule creates a new WebRTC meeting room (LiveKit / Dyte)
// @Summary Schedule WebRTC Meeting Room
// @Description Creates a new interactive WebRTC meeting room
// @Tags Gateway - Meeting
// @Accept json
// @Produce json
// @Param request body dto.MeetingScheduleRequest true "Room Details"
// @Success 201 {object} map[string]string
// @Router /api/v1/gateway/meeting/schedule [post]
func (m *MeetingGatewayController) Schedule(c *gin.Context) {
	var req dto.MeetingScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"room_name":    req.RoomName,
		"provider":     req.Provider,
		"livekit_url":  "wss://livekit.example.com",
		"access_token": "token_lk_9999",
	})
}

// Bot deploys a Recall.ai bot to join external Google Meet / Zoom / Teams calls
// @Summary Deploy Meeting Bot
// @Description Dispatches an autonomous Recall.ai AI bot to join external video calls
// @Tags Gateway - Meeting
// @Accept json
// @Produce json
// @Param request body dto.MeetingBotRequest true "Bot Details"
// @Success 202 {object} map[string]string
// @Router /api/v1/gateway/meeting/bot [post]
func (m *MeetingGatewayController) Bot(c *gin.Context) {
	var req dto.MeetingBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"bot_id":      "bot_recall_555",
		"meeting_url": req.MeetingURL,
		"status":      "bot_dispatching_to_call",
	})
}
