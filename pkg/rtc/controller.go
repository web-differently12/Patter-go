package rtc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
)

type RecallController struct {
	recallService RecallAIService
}

func NewRecallController(svc RecallAIService) *RecallController {
	return &RecallController{recallService: svc}
}

// CreateMeetingProfile godoc
// @Summary      Créer un Profil de Visioconférence (Diarisation, Traduction, Sentiment, Résumé)
// @Tags         Gateway - Meeting Bots (RTC)
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body MeetingProfile true "Configuration du Profil Visio"
// @Success      200 {object} core.APIResponse{data=MeetingProfile}
// @Router       /api/v1/gateway/rtc/profiles [post]
func (ctrl *RecallController) CreateMeetingProfile(c *gin.Context) {
	var req MeetingProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.recallService.CreateMeetingProfile(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// ListMeetingProfiles godoc
// @Summary      Lister les Profils de Visioconférence du Tenant
// @Tags         Gateway - Meeting Bots (RTC)
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]MeetingProfile}
// @Router       /api/v1/gateway/rtc/profiles [get]
func (ctrl *RecallController) ListMeetingProfiles(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.recallService.ListMeetingProfiles(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// CreateBot godoc
// @Summary      Déployer un Bot d'Appel Visio (Zoom, Google Meet, MS Teams, Webex via Recall.ai)
// @Tags         Gateway - Meeting Bots (RTC)
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body CreateMeetingBotRequest true "Configuration du Bot"
// @Success      200 {object} core.APIResponse{data=MeetingBotResponse}
// @Router       /api/v1/gateway/rtc/bot [post]
func (ctrl *RecallController) CreateBot(c *gin.Context) {
	var req CreateMeetingBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.recallService.CreateMeetingBot(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// ListBots godoc
// @Summary      Lister les Bots de visioconférence actifs du Tenant
// @Tags         Gateway - Meeting Bots (RTC)
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]MeetingBotResponse}
// @Router       /api/v1/gateway/rtc/bots [get]
func (ctrl *RecallController) ListBots(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.recallService.ListBots(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// LeaveMeeting godoc
// @Summary      Demander au Bot de quitter la visioconférence
// @Tags         Gateway - Meeting Bots (RTC)
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        id path string true "ID du Bot"
// @Success      200 {object} core.APIResponse
// @Router       /api/v1/gateway/rtc/bots/{id}/leave [post]
func (ctrl *RecallController) LeaveMeeting(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	botID := c.Param("id")

	if err := ctrl.recallService.LeaveMeeting(c.Request.Context(), tenantID, botID); err != nil {
		core.Error(c, http.StatusNotFound, err.Error())
		return
	}

	core.Success(c, gin.H{"message": "Bot left meeting successfully"})
}
