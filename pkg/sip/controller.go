package sip

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
)

type SIPController struct {
	sipService SIPPBXService
}

func NewSIPController(svc SIPPBXService) *SIPController {
	return &SIPController{sipService: svc}
}

// RegisterTrunk godoc
// @Summary      Enregistrer un Trunk SIP Universel (OVH, 3CX, FreeSWITCH, Asterisk, Aircall)
// @Description  Bypass Twilio en connectant directement le PBX / Standard Téléphonique du client
// @Tags         Gateway - Direct SIP Trunking PBX
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body SIPTrunkConfig true "Configuration Trunk SIP"
// @Success      201 {object} core.APIResponse{data=SIPTrunkConfig}
// @Router       /api/v1/gateway/sip/trunks [post]
func (ctrl *SIPController) RegisterTrunk(c *gin.Context) {
	var req SIPTrunkConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.sipService.RegisterTrunk(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Created(c, resp)
}

// ListTrunks godoc
// @Summary      Lister les Trunks SIP configurés pour le Tenant
// @Tags         Gateway - Direct SIP Trunking PBX
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]SIPTrunkConfig}
// @Router       /api/v1/gateway/sip/trunks [get]
func (ctrl *SIPController) ListTrunks(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.sipService.ListTrunks(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// InitiateSIPCall godoc
// @Summary      Emettre un appel via Trunk SIP direct sans Twilio
// @Tags         Gateway - Direct SIP Trunking PBX
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body SIPCallRequest true "Paramètres de l'appel SIP"
// @Success      200 {object} core.APIResponse{data=SIPCallResponse}
// @Router       /api/v1/gateway/sip/call [post]
func (ctrl *SIPController) InitiateSIPCall(c *gin.Context) {
	var req SIPCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.sipService.InitiateSIPCall(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}
