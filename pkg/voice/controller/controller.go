package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
	"github.com/lynxflow/patter-go/pkg/voice/service"
)

type VoiceController struct {
	voiceService service.VoiceService
}

func NewVoiceController(svc service.VoiceService) *VoiceController {
	return &VoiceController{voiceService: svc}
}

// InitiateCall godoc
// @Summary      Initier un appel vocal entrant / sortant avec Agent IA
// @Tags         Gateway - Voice
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body service.VoiceCallRequest true "Paramètres de l'appel"
// @Success      200 {object} core.APIResponse{data=service.VoiceCallResponse}
// @Router       /api/v1/gateway/voice/call [post]
func (ctrl *VoiceController) InitiateCall(c *gin.Context) {
	var req service.VoiceCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.voiceService.InitiateCall(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}
