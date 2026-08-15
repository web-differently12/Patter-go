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

// CreateVoiceProfile godoc
// @Summary      Créer un Profil Vocal complet (Voix, Bruit de fond, Sensibilité, Moteur TTS/STT)
// @Tags         Gateway - Voice
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body service.VoiceProfile true "Configuration du Profil Vocal"
// @Success      200 {object} core.APIResponse{data=service.VoiceProfile}
// @Router       /api/v1/gateway/voice/profiles [post]
func (ctrl *VoiceController) CreateVoiceProfile(c *gin.Context) {
	var req service.VoiceProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.voiceService.CreateVoiceProfile(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// ListVoiceProfiles godoc
// @Summary      Lister les Profils Vocaux du Tenant
// @Tags         Gateway - Voice
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]service.VoiceProfile}
// @Router       /api/v1/gateway/voice/profiles [get]
func (ctrl *VoiceController) ListVoiceProfiles(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.voiceService.ListVoiceProfiles(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// InitiateCall godoc
// @Summary      Initier un appel vocal entrant / sortant avec Agent IA
// @Tags         Gateway - Voice
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body service.InitiateCallRequest true "Paramètres de l'appel"
// @Success      200 {object} core.APIResponse{data=service.InitiateCallResponse}
// @Router       /api/v1/gateway/voice/call [post]
func (ctrl *VoiceController) InitiateCall(c *gin.Context) {
	var req service.InitiateCallRequest
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
