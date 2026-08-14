package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
	"github.com/lynxflow/patter-go/pkg/whatsapp/dto"
	"github.com/lynxflow/patter-go/pkg/whatsapp/service"
)

type WhatsAppController struct {
	waService service.WhatsAppService
}

func NewWhatsAppController(svc service.WhatsAppService) *WhatsAppController {
	return &WhatsAppController{waService: svc}
}

// ConnectSession godoc
// @Summary      Connecter une session WhatsApp (Evolution Go)
// @Description  Inscrit une nouvelle session WhatsApp sous la marque blanche du Tenant
// @Tags         Gateway - WhatsApp
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.ConnectWhatsAppRequest true "Information de la session"
// @Success      200 {object} core.APIResponse{data=dto.QRCodeResponse}
// @Router       /api/v1/gateway/whatsapp/connect [post]
func (ctrl *WhatsAppController) ConnectSession(c *gin.Context) {
	var req dto.ConnectWhatsAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.waService.ConnectSession(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// GetQRCode godoc
// @Summary      Obtenir le QR Code d'une session WhatsApp
// @Description  Renvoie le QR Code Base64 pour authentification marque blanche
// @Tags         Gateway - WhatsApp
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        session path string true "Nom de la session"
// @Success      200 {object} core.APIResponse{data=dto.QRCodeResponse}
// @Router       /api/v1/gateway/whatsapp/qrcode [get]
func (ctrl *WhatsAppController) GetQRCode(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	sessionName := c.Query("session")
	if sessionName == "" {
		sessionName = c.Param("session")
	}
	if sessionName == "" {
		sessionName = "default_wa_1"
	}

	resp, err := ctrl.waService.GetQRCode(c.Request.Context(), tenantID, sessionName)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// ListSessions godoc
// @Summary      Lister les sessions WhatsApp du Tenant
// @Tags         Gateway - WhatsApp
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]dto.SessionStatusResponse}
// @Router       /api/v1/gateway/whatsapp/sessions [get]
func (ctrl *WhatsAppController) ListSessions(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.waService.ListSessions(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// SendMessage godoc
// @Summary      Envoyer un message WhatsApp
// @Tags         Gateway - WhatsApp
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.SendMessageRequest true "Contenu du message"
// @Success      200 {object} core.APIResponse{data=dto.SendMessageResponse}
// @Router       /api/v1/gateway/whatsapp/message/send [post]
func (ctrl *WhatsAppController) SendMessage(c *gin.Context) {
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.waService.SendMessage(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}
