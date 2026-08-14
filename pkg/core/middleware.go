package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates API keys and handles multi-tenant White-Label headers
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader(TenantHeader)
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}
		if tenantID == "" {
			tenantID = "default_tenant"
		}

		apiKey := c.GetHeader(ApiKeyHeader)
		if apiKey == "" {
			apiKey = c.GetHeader("Authorization")
		}

		c.Set(TenantIDKey, tenantID)
		c.Next()
	}
}

// CORSMiddleware provides cross-origin capabilities for white-label dashboards
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-ID, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
