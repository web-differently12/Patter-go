package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/brain"
	"github.com/lynxflow/patter-go/pkg/brain/service"
	"github.com/lynxflow/patter-go/pkg/calendar"
	"github.com/lynxflow/patter-go/pkg/core"
	"github.com/lynxflow/patter-go/pkg/rag"
)

type BrainController struct {
	brainService service.BrainService
}

func NewBrainController(svc service.BrainService) *BrainController {
	return &BrainController{brainService: svc}
}

// QueryBrain godoc
// @Summary      Interroger le Cerveau IA
// @Tags         Gateway - Brain RAG & Analytics
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body service.QueryBrainRequest true "Requête IA"
// @Success      200 {object} core.APIResponse{data=service.QueryBrainResponse}
// @Router       /api/v1/gateway/brain/query [post]
func (ctrl *BrainController) QueryBrain(c *gin.Context) {
	var req service.QueryBrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.brainService.QueryBrain(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// EnhancePrompt godoc
// @Summary      Générer ou Améliorer le Prompt Système d'un Agent IA (Copilot)
// @Description  Optimise les consignes brutes d'un utilisateur en un prompt système structuré prêt pour les appels vocaux et réunions
// @Tags         Gateway - Brain RAG & Analytics
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body service.EnhancePromptRequest true "Prompt brut"
// @Success      200 {object} core.APIResponse{data=service.EnhancePromptResponse}
// @Router       /api/v1/gateway/brain/prompt/enhance [post]
func (ctrl *BrainController) EnhancePrompt(c *gin.Context) {
	var req service.EnhancePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.brainService.EnhanceSystemPrompt(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

type RAGSearchRequest struct {
	CallSID string                 `json:"call_sid"`
	Query   string                 `json:"query" binding:"required"`
	Config  rag.KnowledgeBaseConfig `json:"config"`
}

// SearchRAG godoc
// @Summary      Recherche RAG & Décision d'Escalade
// @Tags         Gateway - Brain RAG & Analytics
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body RAGSearchRequest true "Recherche Vectorielle"
// @Success      200 {object} core.APIResponse{data=rag.SearchResult}
// @Router       /api/v1/gateway/brain/rag/search [post]
func (ctrl *BrainController) SearchRAG(c *gin.Context) {
	var req RAGSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.brainService.SearchRAG(c.Request.Context(), tenantID, req.CallSID, req.Query, req.Config)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

type HumanTransferRequest struct {
	CallSID     string `json:"call_sid" binding:"required"`
	RawArgs     string `json:"raw_args" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

// HumanTransfer godoc
// @Summary      Compétence Système: Transfert Humain avec Consentement
// @Tags         Gateway - Brain RAG & Analytics
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body HumanTransferRequest true "Demande de Transfert"
// @Success      200 {object} core.APIResponse{data=string}
// @Router       /api/v1/gateway/brain/transfer [post]
func (ctrl *BrainController) HumanTransfer(c *gin.Context) {
	var req HumanTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.brainService.ExecuteHumanTransfer(c.Request.Context(), tenantID, req.CallSID, req.RawArgs, req.Destination)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

type CalendarRequest struct {
	RawArgs string                 `json:"raw_args" binding:"required"`
	Config  calendar.CalendarConfig `json:"config"`
}

// CalendarAvailability godoc
// @Summary      Vérifier les disponibilités du calendrier
// @Tags         Gateway - Brain RAG & Analytics
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body CalendarRequest true "Paramètres Calendrier"
// @Success      200 {object} core.APIResponse{data=string}
// @Router       /api/v1/gateway/brain/calendar/availability [post]
func (ctrl *BrainController) CalendarAvailability(c *gin.Context) {
	var req CalendarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.brainService.CheckCalendarAvailability(c.Request.Context(), tenantID, req.RawArgs, req.Config)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// CalendarBook godoc
// @Summary      Réserver un créneau dans le calendrier
// @Tags         Gateway - Brain RAG & Analytics
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body CalendarRequest true "Paramètres Réservation"
// @Success      200 {object} core.APIResponse{data=string}
// @Router       /api/v1/gateway/brain/calendar/book [post]
func (ctrl *BrainController) CalendarBook(c *gin.Context) {
	var req CalendarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.brainService.BookCalendarAppointment(c.Request.Context(), tenantID, req.RawArgs, req.Config)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

type LeMURRequest struct {
	TranscriptID string                  `json:"transcript_id" binding:"required"`
	Utterances   []brain.SpeakerUtterance `json:"utterances"`
}

// ProcessLeMUR godoc
// @Summary      Analyse Post-Appel AssemblyAI LeMUR v3
// @Tags         Gateway - Brain RAG & Analytics
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body LeMURRequest true "Transcription brute"
// @Success      200 {object} core.APIResponse{data=brain.LeMURResponse}
// @Router       /api/v1/gateway/brain/lemur/process [post]
func (ctrl *BrainController) ProcessLeMUR(c *gin.Context) {
	var req LeMURRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.brainService.ProcessPostCallLeMUR(c.Request.Context(), tenantID, req.TranscriptID, req.Utterances)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}
