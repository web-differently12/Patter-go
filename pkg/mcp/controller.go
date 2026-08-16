package mcp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
)

type MCPController struct {
	mcpService MCPService
	hub        *TenantIntegrationHub
}

func NewMCPController(svc MCPService) *MCPController {
	return &MCPController{
		mcpService: svc,
		hub:        NewTenantIntegrationHub(svc),
	}
}

// ListTemplates godoc
// @Summary      Lister le catalogue des modèles d'intégration (Nango Template Catalogue)
// @Tags         Gateway - Model Context Protocol (MCP)
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=[]IntegrationTemplate}
// @Router       /api/v1/gateway/integrations/templates [get]
func (ctrl *MCPController) ListTemplates(c *gin.Context) {
	resp := ctrl.hub.ListTemplates()
	core.Success(c, resp)
}

// DeployTemplate godoc
// @Summary      Déployer un modèle d'intégration en 1-Click (HubSpot, Salesforce, Google Workspace, Slack)
// @Tags         Gateway - Model Context Protocol (MCP)
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body DeployTemplateRequest true "Information du modèle"
// @Success      200 {object} core.APIResponse{data=DeployTemplateResponse}
// @Router       /api/v1/gateway/integrations/templates/deploy [post]
func (ctrl *MCPController) DeployTemplate(c *gin.Context) {
	var req DeployTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.hub.DeployTemplate(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// GetN8NCommunityNodeSchema godoc
// @Summary      Obtenir le schéma du nœud communautaire n8n / Make.com 1-Click
// @Description  Génère automatiquement les spécifications de nœuds n8n pour les serveurs MCP et webhooks du Tenant
// @Tags         Gateway - Model Context Protocol (MCP)
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Success      200 {object} core.APIResponse{data=N8NCommunityNodeSchema}
// @Router       /api/v1/gateway/integrations/n8n/node-schema [get]
func (ctrl *MCPController) GetN8NCommunityNodeSchema(c *gin.Context) {
	tenantID := core.GetTenantID(c)
	resp, err := ctrl.hub.GenerateN8NNodeSchema(c.Request.Context(), tenantID)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
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
