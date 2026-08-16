package core

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	TenantIDKey  = "tenant_id"
	ApiKeyHeader = "X-API-Key"
	TenantHeader = "X-Tenant-ID"
)

type contextKey string

const tenantCtxKey contextKey = "tenant_id"

// WithTenantID injects tenant_id into standard context.Context
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantCtxKey, tenantID)
}

// TenantIDFromContext retrieves tenant_id from standard context.Context
func TenantIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(tenantCtxKey).(string); ok {
		return val
	}
	return ""
}

// GetTenantID retrieves tenant_id from Gin context
func GetTenantID(c *gin.Context) string {
	if val, exists := c.Get(TenantIDKey); exists {
		if id, ok := val.(string); ok {
			return id
		}
	}
	header := c.GetHeader(TenantHeader)
	if header != "" {
		return header
	}
	return "default"
}
