package mcp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type N8NNodeProperty struct {
	DisplayName string      `json:"displayName"`
	Name        string      `json:"name"`
	Type        string      `json:"type"` // "string", "options", "boolean", "collection"
	Default     interface{} `json:"default"`
	Description string      `json:"description,omitempty"`
}

type N8NCommunityNodeSchema struct {
	NodeName        string            `json:"nodeName"`
	DisplayName     string            `json:"displayName"`
	Icon            string            `json:"icon"`
	Description     string            `json:"description"`
	Version         int               `json:"version"`
	Defaults        map[string]string `json:"defaults"`
	Inputs          []string          `json:"inputs"`  // ["main"]
	Outputs         []string          `json:"outputs"` // ["main"]
	WebhookTrigger  string            `json:"webhookTriggerUrl"`
	CredentialsType string            `json:"credentialsType"` // "patterApi"
	Properties      []N8NNodeProperty `json:"properties"`
}

type IntegrationTemplate struct {
	TemplateID    string   `json:"template_id"`
	Title         string   `json:"title"`
	Category      string   `json:"category"` // "CRM", "CALENDAR", "MESSAGING", "ERP_HRIS"
	Description   string   `json:"description"`
	Provider      string   `json:"provider"` // "HubSpot", "Salesforce", "Google Workspace", "Slack", "Zendesk"
	AuthType      string   `json:"auth_type"` // "OAuth2", "APIKey"
	ExposedTools  []string `json:"exposed_tools"`
	PopularRating float64  `json:"popular_rating"`
}

type DeployTemplateRequest struct {
	TemplateID        string            `json:"template_id" binding:"required"`
	CustomName        string            `json:"custom_name,omitempty"`
	CredentialsConfig map[string]string `json:"credentials_config,omitempty"`
}

type DeployTemplateResponse struct {
	IntegrationID string    `json:"integration_id"`
	ServerID      string    `json:"server_id"`
	Status        string    `json:"status"` // "ACTIVE", "PENDING_OAUTH"
	OAuthURL      string    `json:"oauth_url,omitempty"`
	DeployedAt    time.Time `json:"deployed_at"`
}

type ConnectTenantAPIRequest struct {
	Name        string            `json:"name" binding:"required"` // e.g. "Mon CRM Interne"
	Provider    string            `json:"provider" binding:"required"` // e.g. "HubSpot", "CustomAPI"
	AuthType    string            `json:"auth_type"` // "APIKey", "OAuth2"
	APIKey      string            `json:"api_key,omitempty"`
	BaseURL     string            `json:"base_url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

type ActiveTenantIntegration struct {
	IntegrationID string            `json:"integration_id"`
	TenantID      string            `json:"tenant_id"`
	Name          string            `json:"name"`
	Provider      string            `json:"provider"`
	Status        string            `json:"status"` // "CONNECTED", "DISCONNECTED"
	ServerID      string            `json:"server_id"`
	CreatedAt     time.Time         `json:"created_at"`
}

type TenantIntegrationHub struct {
	mu          sync.RWMutex
	mcpService  MCPService
	templates   []IntegrationTemplate
	activeInts  map[string]*ActiveTenantIntegration
}

func NewTenantIntegrationHub(mcpSvc MCPService) *TenantIntegrationHub {
	templates := []IntegrationTemplate{
		{
			TemplateID:    "tpl_hubspot_crm_sync",
			Title:         "HubSpot CRM Sync & Live Lead Creation",
			Category:      "CRM",
			Description:   "Synchronise vos contacts, opportunités et crée automatiquement des tickets depuis vos appels vocaux et visioconférences",
			Provider:      "HubSpot",
			AuthType:      "OAuth2",
			ExposedTools:  []string{"get_crm_contact", "create_lead", "update_deal_stage", "create_support_ticket"},
			PopularRating: 4.9,
		},
		{
			TemplateID:    "tpl_salesforce_enterprise",
			Title:         "Salesforce Enterprise Account Manager",
			Category:      "CRM",
			Description:   "Met à jour le pipeline commercial Salesforce et enregistre le score BANT post-appel",
			Provider:      "Salesforce",
			AuthType:      "OAuth2",
			ExposedTools:  []string{"query_account", "create_opportunity", "log_call_transcript"},
			PopularRating: 4.8,
		},
		{
			TemplateID:    "tpl_google_workspace_calendar",
			Title:         "Google Workspace & Calendar Auto-Booking",
			Category:      "CALENDAR",
			Description:   "Vérifie vos disponibilités Google Calendar et réserve des créneaux en direct pendant un appel vocal",
			Provider:      "Google Workspace",
			AuthType:      "OAuth2",
			ExposedTools:  []string{"check_free_slots", "book_google_event", "send_invite_email"},
			PopularRating: 5.0,
		},
		{
			TemplateID:    "tpl_slack_notifications",
			Title:         "Slack Live Notifications & Escalation Alerts",
			Category:      "MESSAGING",
			Description:   "Envoie des alertes en direct dans un canal Slack lors d'une escalade vers un conseiller humain ou d'un résumé de réunion",
			Provider:      "Slack",
			AuthType:      "OAuth2",
			ExposedTools:  []string{"post_channel_message", "send_direct_alert"},
			PopularRating: 4.7,
		},
	}

	return &TenantIntegrationHub{
		mcpService: mcpSvc,
		templates:  templates,
		activeInts: make(map[string]*ActiveTenantIntegration),
	}
}

func (h *TenantIntegrationHub) ListTemplates() []IntegrationTemplate {
	return h.templates
}

func (h *TenantIntegrationHub) ConnectTenantAPIIntegration(ctx context.Context, tenantID string, req ConnectTenantAPIRequest) (*ActiveTenantIntegration, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	integrationID := "int_api_" + uuid.New().String()[:8]
	serverID := "mcp_srv_" + uuid.New().String()[:8]

	headers := req.Headers
	if headers == nil {
		headers = make(map[string]string)
	}
	if req.APIKey != "" {
		headers["Authorization"] = "Bearer " + req.APIKey
	}

	baseURL := req.BaseURL
	if baseURL == "" {
		baseURL = "https://api.patter.ai/v1/bridge/api/" + integrationID
	}

	_, err := h.mcpService.RegisterServer(ctx, tenantID, MCPServerConfig{
		ServerID:            serverID,
		TenantID:            tenantID,
		Name:                req.Name,
		URL:                 baseURL,
		Transport:           TransportPatterBridge,
		PatterIntegrationID: integrationID,
		AuthHeader:          headers,
		Status:              "CONNECTED",
	})
	if err != nil {
		return nil, err
	}

	active := &ActiveTenantIntegration{
		IntegrationID: integrationID,
		TenantID:      tenantID,
		Name:          req.Name,
		Provider:      req.Provider,
		Status:        "CONNECTED",
		ServerID:      serverID,
		CreatedAt:     time.Now(),
	}

	h.activeInts[tenantID+":"+integrationID] = active
	return active, nil
}

func (h *TenantIntegrationHub) ListActiveTenantIntegrations(ctx context.Context, tenantID string) ([]*ActiveTenantIntegration, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []*ActiveTenantIntegration
	for _, item := range h.activeInts {
		if item.TenantID == tenantID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (h *TenantIntegrationHub) DeployTemplate(ctx context.Context, tenantID string, req DeployTemplateRequest) (*DeployTemplateResponse, error) {
	integrationID := "int_" + uuid.New().String()[:8]
	serverID := "mcp_tpl_" + uuid.New().String()[:8]

	_, err := h.mcpService.RegisterServer(ctx, tenantID, MCPServerConfig{
		ServerID:            serverID,
		TenantID:            tenantID,
		Name:                req.CustomName,
		URL:                 "https://api.patter.ai/v1/bridge/template/" + req.TemplateID,
		Transport:           TransportPatterBridge,
		PatterIntegrationID: integrationID,
		Status:              "CONNECTED",
	})
	if err != nil {
		return nil, err
	}

	return &DeployTemplateResponse{
		IntegrationID: integrationID,
		ServerID:      serverID,
		Status:        "ACTIVE",
		OAuthURL:      fmt.Sprintf("https://auth.patter.ai/v1/oauth/connect/%s?tenant_id=%s", req.TemplateID, tenantID),
		DeployedAt:    time.Now(),
	}, nil
}

func (h *TenantIntegrationHub) GenerateN8NNodeSchema(ctx context.Context, tenantID string) (*N8NCommunityNodeSchema, error) {
	servers, err := h.mcpService.ListServers(ctx, tenantID)
	if err != nil {
		servers = []*MCPServerConfig{}
	}

	var mcpServerOptions []string
	for _, srv := range servers {
		mcpServerOptions = append(mcpServerOptions, srv.Name)
	}

	return &N8NCommunityNodeSchema{
		NodeName:        "n8n-nodes-patter-gateway",
		DisplayName:     "Patter Omnichannel Gateway",
		Icon:            "file:patter.svg",
		Description:     "Noeud communautaire 1-Click n8n / Make.com pour Patter Core Engine Gateway",
		Version:         1,
		Defaults:        map[string]string{"name": "Patter Gateway"},
		Inputs:          []string{"main"},
		Outputs:         []string{"main"},
		WebhookTrigger:  fmt.Sprintf("https://api.patter.ai/api/v1/gateway/webhooks/n8n/%s", tenantID),
		CredentialsType: "patterApiHeader",
		Properties: []N8NNodeProperty{
			{
				DisplayName: "Ressource Gateway",
				Name:        "resource",
				Type:        "options",
				Default:     "voiceCall",
				Description: "Sélectionnez l'action Patter à exécuter dans n8n",
			},
			{
				DisplayName: "Numéro Destinataire (E.164)",
				Name:        "recipient",
				Type:        "string",
				Default:     "+33612345678",
				Description: "Numéro de téléphone du client",
			},
			{
				DisplayName: "Canal de Message",
				Name:        "channel",
				Type:        "options",
				Default:     "WHATSAPP",
				Description: "Protocole d'envoi (WHATSAPP, SMS, RCS)",
			},
			{
				DisplayName: "Serveurs MCP Connectés",
				Name:        "mcpServer",
				Type:        "options",
				Default:     "Patter Unified Integrations Bridge",
				Description: "Serveur MCP exécuté de manière transparente",
			},
		},
	}, nil
}
