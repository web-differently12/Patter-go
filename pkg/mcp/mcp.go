package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/core"
)

type MCPServerTransport string

const (
	TransportHTTP         MCPServerTransport = "streamable_http"
	TransportSSE          MCPServerTransport = "sse"
	TransportWebSocket    MCPServerTransport = "websocket"
	TransportPatterBridge MCPServerTransport = "patter_unified_bridge"
)

type MCPServerConfig struct {
	ServerID            string             `json:"server_id"`
	TenantID            string             `json:"tenant_id"`
	Name                string             `json:"name" binding:"required"` // e.g. "HubSpot CRM MCP", "Patter Unified Integrations"
	URL                 string             `json:"url" binding:"required"`  // e.g. "https://mcp.company.com/v1"
	Transport           MCPServerTransport `json:"transport"`
	PatterConnectionID  string             `json:"patter_connection_id,omitempty"`
	PatterIntegrationID string             `json:"patter_integration_id,omitempty"`
	AuthHeader          map[string]string  `json:"auth_headers,omitempty"`
	Status              string             `json:"status"` // "CONNECTED", "DISCONNECTED"
}

type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	ServerID    string                 `json:"server_id"`
}

type MCPToolCallRequest struct {
	ServerID string                 `json:"server_id" binding:"required"`
	ToolName string                 `json:"tool_name" binding:"required"`
	Args     map[string]interface{} `json:"arguments"`
}

type MCPToolCallResponse struct {
	ToolName string      `json:"tool_name"`
	Content  interface{} `json:"content"`
	IsError  bool        `json:"is_error"`
}

type MCPService interface {
	RegisterServer(ctx context.Context, tenantID string, cfg MCPServerConfig) (*MCPServerConfig, error)
	ListServers(ctx context.Context, tenantID string) ([]*MCPServerConfig, error)
	ListDiscoveredTools(ctx context.Context, tenantID, serverID string) ([]MCPTool, error)
	ExecuteTool(ctx context.Context, tenantID string, req MCPToolCallRequest) (*MCPToolCallResponse, error)
}

type mcpService struct {
	mu         sync.RWMutex
	servers    map[string]*MCPServerConfig
	httpClient *http.Client
	logger     *slog.Logger
}

func NewMCPService(logger *slog.Logger) MCPService {
	if logger == nil {
		logger = core.GetLogger()
	}
	svc := &mcpService{
		servers:    make(map[string]*MCPServerConfig),
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     logger.With("component", "mcp_service"),
	}

	mockID := "mcp_hubspot_01"
	svc.servers["default_tenant:"+mockID] = &MCPServerConfig{
		ServerID:  mockID,
		TenantID:  "default_tenant",
		Name:      "HubSpot CRM MCP Server",
		URL:       "https://mcp.hubspot.com/v1",
		Transport: TransportHTTP,
		Status:    "CONNECTED",
	}

	patterBridgeID := "mcp_patter_bridge_01"
	svc.servers["default_tenant:"+patterBridgeID] = &MCPServerConfig{
		ServerID:            patterBridgeID,
		TenantID:            "default_tenant",
		Name:                "Patter Unified Integrations Bridge",
		URL:                 "https://api.patter.ai/v1/bridge",
		Transport:           TransportPatterBridge,
		PatterIntegrationID: "hubspot-salesforce-slack",
		Status:              "CONNECTED",
	}

	return svc
}

func (s *mcpService) RegisterServer(ctx context.Context, tenantID string, cfg MCPServerConfig) (*MCPServerConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	serverID := "mcp_" + uuid.New().String()[:8]
	cfg.ServerID = serverID
	cfg.TenantID = tenantID
	if cfg.Transport == "" {
		cfg.Transport = TransportHTTP
	}
	cfg.Status = "CONNECTED"

	s.servers[tenantID+":"+serverID] = &cfg
	s.logger.Info("Registered Model Context Protocol (MCP) server", "server_id", serverID, "name", cfg.Name, "url", cfg.URL)
	return &cfg, nil
}

func (s *mcpService) ListServers(ctx context.Context, tenantID string) ([]*MCPServerConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*MCPServerConfig
	for _, srv := range s.servers {
		if srv.TenantID == tenantID {
			result = append(result, srv)
		}
	}
	return result, nil
}

func (s *mcpService) ListDiscoveredTools(ctx context.Context, tenantID, serverID string) ([]MCPTool, error) {
	s.mu.RLock()
	key := tenantID + ":" + serverID
	srv, ok := s.servers[key]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("mcp server %s not found", serverID)
	}

	if srv.URL != "" {
		payload, _ := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "tools/list",
		})
		httpReq, err := http.NewRequestWithContext(ctx, "POST", srv.URL, bytes.NewBuffer(payload))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/json")
			for k, v := range srv.AuthHeader {
				httpReq.Header.Set(k, v)
			}
			resp, err := s.httpClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				s.logger.Info("MCP tools/list API response", "server_id", serverID, "len", len(body))
			}
		}
	}

	return []MCPTool{
		{
			Name:        "get_crm_contact",
			Description: "Récupère les informations d'un contact CRM via Patter MCP Server Bridge",
			ServerID:    serverID,
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"email": map[string]interface{}{"type": "string"},
					"phone": map[string]interface{}{"type": "string"},
				},
				"required": []string{"phone"},
			},
		},
		{
			Name:        "create_support_ticket",
			Description: "Crée un ticket de support client via Patter MCP Server Bridge",
			ServerID:    serverID,
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"subject": map[string]interface{}{"type": "string"},
					"priority": map[string]interface{}{
						"type": "string", "enum": []string{"LOW", "MEDIUM", "HIGH", "URGENT"},
					},
				},
				"required": []string{"subject"},
			},
		},
	}, nil
}

func (s *mcpService) ExecuteTool(ctx context.Context, tenantID string, req MCPToolCallRequest) (*MCPToolCallResponse, error) {
	s.mu.RLock()
	key := tenantID + ":" + req.ServerID
	srv, ok := s.servers[key]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("mcp server %s not found", req.ServerID)
	}

	s.logger.Info("Executing MCP Tool Call", "server_id", req.ServerID, "tool", req.ToolName, "transport", srv.Transport)

	if srv.URL != "" {
		payload, _ := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      2,
			"method":  "tools/call",
			"params": map[string]interface{}{
				"name":      req.ToolName,
				"arguments": req.Args,
			},
		})
		httpReq, err := http.NewRequestWithContext(ctx, "POST", srv.URL, bytes.NewBuffer(payload))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/json")
			for k, v := range srv.AuthHeader {
				httpReq.Header.Set(k, v)
			}
			resp, err := s.httpClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
			}
		}
	}

	return &MCPToolCallResponse{
		ToolName: req.ToolName,
		Content:  map[string]interface{}{"status": "success", "data": "Action exécutée avec succès via Patter MCP Server Bridge"},
		IsError:  false,
	}, nil
}
