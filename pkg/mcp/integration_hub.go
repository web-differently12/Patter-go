package mcp

import (
	"context"
	"fmt"
)

type N8NNodeProperty struct {
	DisplayName string   `json:"displayName"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // "string", "options", "boolean", "collection"
	Default     interface{} `json:"default"`
	Description string   `json:"description,omitempty"`
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

type TenantIntegrationHub struct {
	mcpService MCPService
}

func NewTenantIntegrationHub(mcpSvc MCPService) *TenantIntegrationHub {
	return &TenantIntegrationHub{mcpService: mcpSvc}
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
