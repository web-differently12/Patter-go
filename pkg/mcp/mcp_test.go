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
