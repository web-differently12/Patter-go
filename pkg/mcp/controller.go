package mcp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
)

type MCPController struct {
	mcpService MCPService
}

func NewMCPController(svc MCPService) *MCPController {
	return &MCPController{mcpService: svc}
}

// RegisterServer godoc
// @Summary      Connecter un serveur MCP (Model Context Protocol)
// @Description  Intègre des outils externes (HubSpot, Postgres, GitHub, Slack) via MCP
// @Tags         Gateway - Model Context Protocol (MCP)
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body MCPServerConfig true "Configuration Serveur MCP"
// @Success      201 {object} core.APIResponse{data=MCPServerConfig}
// @Router       /api/v1/gateway/mcp/servers [post]
func (ctrl *MCPController) RegisterServer(c *gin.Context) {
	var req MCPServerConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.mcpService.RegisterServer(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Created(c, resp)
}

// ListServers godoc
// @Summary      Lister les serveurs MCP connectés au Tenant
// @Tags         Gateway - Model Context Protocol (MCP)
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]MCPServerConfig}
// @Router       /api/v1/gateway/mcp/servers [get]
func (ctrl *MCPController) ListServers(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.mcpService.ListServers(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// ListTools godoc
// @Summary      Découvrir les outils exposés par un serveur MCP (tools/list)
// @Tags         Gateway - Model Context Protocol (MCP)
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        id path string true "ID du Serveur MCP"
// @Success      200 {object} core.APIResponse{data=[]MCPTool}
// @Router       /api/v1/gateway/mcp/servers/{id}/tools [get]
func (ctrl *MCPController) ListTools(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	serverID := c.Param("id")

	resp, err := ctrl.mcpService.ListDiscoveredTools(c.Request.Context(), tenantID, serverID)
	if err != nil {
		core.Error(c, http.StatusNotFound, err.Error())
		return
	}

	core.Success(c, resp)
}

// ExecuteTool godoc
// @Summary      Exécuter un outil MCP (tools/call)
// @Tags         Gateway - Model Context Protocol (MCP)
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body MCPToolCallRequest true "Arguments de l'outil"
// @Success      200 {object} core.APIResponse{data=MCPToolCallResponse}
// @Router       /api/v1/gateway/mcp/tools/execute [post]
func (ctrl *MCPController) ExecuteTool(c *gin.Context) {
	var req MCPToolCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.mcpService.ExecuteTool(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}
