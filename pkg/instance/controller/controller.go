package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
	"github.com/lynxflow/patter-go/pkg/instance/dto"
	"github.com/lynxflow/patter-go/pkg/instance/service"
)

type InstanceController struct {
	instanceService service.InstanceService
}

func NewInstanceController(svc service.InstanceService) *InstanceController {
	return &InstanceController{instanceService: svc}
}

// CreateInstance godoc
// @Summary      Créer une instance White-Label
// @Description  Initialise un nouvel espace d'isolation pour le tenant
// @Tags         Gateway - Instances
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body dto.CreateInstanceRequest true "Configuration de l'instance"
// @Success      201 {object} core.APIResponse{data=dto.InstanceResponse}
// @Failure      400 {object} core.APIResponse
// @Router       /api/v1/gateway/instances [post]
func (ctrl *InstanceController) CreateInstance(c *gin.Context) {
	var req dto.CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.instanceService.CreateInstance(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Created(c, resp)
}

// ListInstances godoc
// @Summary      Lister les instances du Tenant
// @Tags         Gateway - Instances
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]dto.InstanceResponse}
// @Router       /api/v1/gateway/instances [get]
func (ctrl *InstanceController) ListInstances(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.instanceService.ListInstances(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// GetInstance godoc
// @Summary      Obtenir les détails d'une instance
// @Tags         Gateway - Instances
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        id path string true "ID de l'instance"
// @Success      200 {object} core.APIResponse{data=dto.InstanceResponse}
// @Router       /api/v1/gateway/instances/{id} [get]
func (ctrl *InstanceController) GetInstance(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	instID := c.Param("id")

	resp, err := ctrl.instanceService.GetInstance(c.Request.Context(), tenantID, instID)
	if err != nil {
		core.Error(c, http.StatusNotFound, err.Error())
		return
	}

	core.Success(c, resp)
}

// DeleteInstance godoc
// @Summary      Supprimer une instance
// @Tags         Gateway - Instances
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        id path string true "ID de l'instance"
// @Success      200 {object} core.APIResponse
// @Router       /api/v1/gateway/instances/{id} [delete]
func (ctrl *InstanceController) DeleteInstance(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	instID := c.Param("id")

	if err := ctrl.instanceService.DeleteInstance(c.Request.Context(), tenantID, instID); err != nil {
		core.Error(c, http.StatusNotFound, err.Error())
		return
	}

	core.Success(c, gin.H{"message": "Instance deleted successfully"})
}
