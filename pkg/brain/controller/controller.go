package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/brain/service"
	"github.com/lynxflow/patter-go/pkg/core"
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
