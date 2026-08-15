package gateway

import (
	"net/http"

	"vocal-engine/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

type AvatarGatewayController struct{}

func NewAvatarGatewayController() *AvatarGatewayController {
	return &AvatarGatewayController{}
}

// Live generates a WebRTC live session token for interactive avatar conversation
// @Summary Start Live WebRTC Avatar Session
// @Description Generates a WebRTC join token for interactive real-time avatar streams
// @Tags Gateway - Avatar
// @Accept json
// @Produce json
// @Param request body dto.AvatarLiveRequest true "Session Specs"
// @Success 200 {object} map[string]string
// @Router /api/v1/gateway/avatar/live [post]
func (a *AvatarGatewayController) Live(c *gin.Context) {
	var req dto.AvatarLiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": "avatar_sess_12345",
		"room_name":  "room_avatar_1",
		"join_token": "token_webrtc_secret",
	})
}

// Offline requests on-the-fly MP4 raw video rendering
// @Summary Request Offline Raw MP4 Video Generation
// @Description Requests on-the-fly MP4 video generation metered strictly per second
// @Tags Gateway - Avatar
// @Accept json
// @Produce json
// @Param request body dto.AvatarOfflineRequest true "Video Specs"
// @Success 202 {object} map[string]string
// @Router /api/v1/gateway/avatar/offline [post]
func (a *AvatarGatewayController) Offline(c *gin.Context) {
	var req dto.AvatarOfflineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":       "avatar_render_job_789",
		"status":       "rendering_mp4_on_the_fly",
		"content_type": "video/mp4",
	})
}
