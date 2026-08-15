package mcp_test

import (
	"context"
	"testing"

	"github.com/lynxflow/patter-go/pkg/mcp"
)

func TestMCPServiceLifecycle(t *testing.T) {
	svc := mcp.NewMCPService(nil)

	cfg := mcp.MCPServerConfig{
		Name:      "Test Postgres MCP",
		URL:       "https://mcp.postgres.local/v1",
		Transport: mcp.TransportHTTP,
	}

	srv, err := svc.RegisterServer(context.Background(), "tenant_test", cfg)
	if err != nil {
		t.Fatalf("Unexpected error registering MCP server: %v", err)
	}

	if srv.ServerID == "" {
		t.Errorf("Expected valid ServerID")
	}

	tools, err := svc.ListDiscoveredTools(context.Background(), "tenant_test", srv.ServerID)
	if err != nil {
		t.Fatalf("Unexpected error listing tools: %v", err)
	}

	if len(tools) == 0 {
		t.Errorf("Expected discovered MCP tools, got 0")
	}

	execReq := mcp.MCPToolCallRequest{
		ServerID: srv.ServerID,
		ToolName: "get_crm_contact",
		Args:     map[string]interface{}{"phone": "+33612345678"},
	}

	resp, err := svc.ExecuteTool(context.Background(), "tenant_test", execReq)
	if err != nil {
		t.Fatalf("Unexpected error executing MCP tool: %v", err)
	}

	if resp.IsError {
		t.Errorf("Expected IsError == false")
	}
}

func TestN8NCommunityNodeSchema(t *testing.T) {
	svc := mcp.NewMCPService(nil)
	hub := mcp.NewTenantIntegrationHub(svc)

	schema, err := hub.GenerateN8NNodeSchema(context.Background(), "tenant_test")
	if err != nil {
		t.Fatalf("Unexpected error generating n8n community node schema: %v", err)
	}

	if schema.NodeName != "n8n-nodes-patter-gateway" {
		t.Errorf("Expected n8n node name n8n-nodes-patter-gateway, got %s", schema.NodeName)
	}
	if len(schema.Properties) == 0 {
		t.Errorf("Expected non-empty properties for 1-click n8n node")
	}
}
