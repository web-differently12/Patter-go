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

// Schedule creates a new WebRTC meeting room
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
		"tier":         req.Tier,
		"join_url":     "wss://meet.example.com",
		"access_token": "token_meet_9999",
	})
}

// Bot deploys an AI assistant bot to join external calls
// @Summary Deploy Meeting Bot
// @Description Dispatches an autonomous AI bot to join external video calls
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
		"bot_id":      "bot_ai_555",
		"meeting_url": req.MeetingURL,
		"status":      "bot_dispatching_to_call",
	})
}
