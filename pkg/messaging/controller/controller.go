package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
	"github.com/lynxflow/patter-go/pkg/messaging/service"
)

type MessagingController struct {
	messagingService service.MessagingService
}

func NewMessagingController(svc service.MessagingService) *MessagingController {
	return &MessagingController{messagingService: svc}
}

// SendSMS godoc
// @Summary      Envoyer un SMS / MMS
// @Tags         Gateway - Messaging
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body service.SendSMSRequest true "Paramètres du SMS"
// @Success      200 {object} core.APIResponse{data=service.SendSMSResponse}
// @Router       /api/v1/gateway/messaging/sms [post]
func (ctrl *MessagingController) SendSMS(c *gin.Context) {
	var req service.SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.messagingService.SendSMS(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// SendRCS godoc
// @Summary      Envoyer un message RCS (Rich Communication Services)
// @Description  Message enrichi avec cartes, carrousels, boutons interactifs et fallback SMS automatique
// @Tags         Gateway - Messaging
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body service.SendRCSRequest true "Paramètres RCS"
// @Success      200 {object} core.APIResponse{data=service.SendRCSResponse}
// @Router       /api/v1/gateway/messaging/rcs [post]
func (ctrl *MessagingController) SendRCS(c *gin.Context) {
	var req service.SendRCSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.messagingService.SendRCS(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}
