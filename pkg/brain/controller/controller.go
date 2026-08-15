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
// @Summary      Interroger le cerveau IA RAG & LeMUR v3
// @Tags         Gateway - Brain
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

type RAGSearchRequest struct {
	CallSID string                  `json:"call_sid,omitempty"`
	Query   string                  `json:"query" binding:"required"`
	Config  rag.KnowledgeBaseConfig `json:"config"`
}

// SearchRAG godoc
// @Summary      Recherche RAG Unifiée avec règle d'escalade automatique
// @Tags         Gateway - Brain
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body RAGSearchRequest true "Paramètres RAG"
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

type TransferRequest struct {
	CallSID     string `json:"call_sid" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	RawArgs     string `json:"raw_args" binding:"required"`
}

// HumanTransfer godoc
// @Summary      Invoquer le skill de transfert vers un humain avec consentement oral
// @Tags         Gateway - Brain
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body TransferRequest true "Arguments de transfert"
// @Success      200 {object} core.APIResponse
// @Router       /api/v1/gateway/brain/transfer [post]
func (ctrl *BrainController) HumanTransfer(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	result, err := ctrl.brainService.ExecuteHumanTransfer(c.Request.Context(), tenantID, req.CallSID, req.RawArgs, req.Destination)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, gin.H{"status": result})
}

type CalendarRequest struct {
	RawArgs string                  `json:"raw_args" binding:"required"`
	Config  calendar.CalendarConfig `json:"config"`
}

// CalendarAvailability godoc
// @Summary      Vérifier la disponibilité de l'agenda (Google Calendar, Outlook, Cal.com, Calendly)
// @Tags         Gateway - Brain
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body CalendarRequest true "Arguments agenda"
// @Success      200 {object} core.APIResponse
// @Router       /api/v1/gateway/brain/calendar/availability [post]
func (ctrl *BrainController) CalendarAvailability(c *gin.Context) {
	var req CalendarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	result, err := ctrl.brainService.CheckCalendarAvailability(c.Request.Context(), tenantID, req.RawArgs, req.Config)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, gin.H{"result": result})
}

// CalendarBook godoc
// @Summary      Réserver un créneau dans l'agenda
// @Tags         Gateway - Brain
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body CalendarRequest true "Arguments réservation"
// @Success      200 {object} core.APIResponse
// @Router       /api/v1/gateway/brain/calendar/book [post]
func (ctrl *BrainController) CalendarBook(c *gin.Context) {
	var req CalendarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	result, err := ctrl.brainService.BookCalendarAppointment(c.Request.Context(), tenantID, req.RawArgs, req.Config)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, gin.H{"result": result})
}

type LeMURRequest struct {
	TranscriptID string                  `json:"transcript_id" binding:"required"`
	Utterances   []brain.SpeakerUtterance `json:"utterances,omitempty"`
}

// ProcessLeMUR godoc
// @Summary      Exécuter l'analyse Post-Call LeMUR v3 (BANT, Action Items, Talk-Ratio & Sentiment)
// @Tags         Gateway - Brain
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body LeMURRequest true "Paramètres de l'appel"
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
