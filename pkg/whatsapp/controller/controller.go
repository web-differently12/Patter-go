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

// CheckNumberExists godoc
// @Summary      Vérifier si un numéro existe sur WhatsApp (JID)
// @Description  Normalise et vérifie la présence d'un numéro de téléphone sur WhatsApp
// @Tags         Gateway - WhatsApp
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.CheckNumberRequest true "Numéro à vérifier"
// @Success      200 {object} core.APIResponse{data=dto.CheckNumberResponse}
// @Router       /api/v1/gateway/whatsapp/check-number [post]
func (ctrl *WhatsAppController) CheckNumberExists(c *gin.Context) {
	var req dto.CheckNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.waService.CheckNumberExists(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
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
// @Summary      Envoyer un message WhatsApp texte
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

// SendMedia godoc
// @Summary      Envoyer un média WhatsApp (Image, Vidéo, Document, Audio, Sticker)
// @Tags         Gateway - WhatsApp
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.SendMediaRequest true "Paramètres média"
// @Success      200 {object} core.APIResponse{data=dto.SendMessageResponse}
// @Router       /api/v1/gateway/whatsapp/media/send [post]
func (ctrl *WhatsAppController) SendMedia(c *gin.Context) {
	var req dto.SendMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.waService.SendMedia(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// SendLocation godoc
// @Summary      Envoyer une localisation GPS WhatsApp
// @Tags         Gateway - WhatsApp
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.SendLocationRequest true "Paramètres localisation"
// @Success      200 {object} core.APIResponse{data=dto.SendMessageResponse}
// @Router       /api/v1/gateway/whatsapp/location/send [post]
func (ctrl *WhatsAppController) SendLocation(c *gin.Context) {
	var req dto.SendLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.waService.SendLocation(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// SendContact godoc
// @Summary      Envoyer une carte de contact VCard WhatsApp
// @Tags         Gateway - WhatsApp
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.SendContactRequest true "Paramètres contact"
// @Success      200 {object} core.APIResponse{data=dto.SendMessageResponse}
// @Router       /api/v1/gateway/whatsapp/contact/send [post]
func (ctrl *WhatsAppController) SendContact(c *gin.Context) {
	var req dto.SendContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.waService.SendContact(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}
