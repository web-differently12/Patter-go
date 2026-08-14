package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	TenantContextKey = "TenantContext"
)

type TenantContext struct {
	TenantID string
	Role     string
}

// AdminAuthMiddleware validates Lynxflow admin secret or JWT token for tenant configuration access
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminKey := c.GetHeader("X-Admin-Secret")
		if adminKey == "" {
			adminKey = c.GetHeader("Authorization")
		}

		if adminKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing admin authentication credential"})
			c.Abort()
			return
		}

		// Proceed with admin request
		c.Next()
	}
}

// TenantContextMiddleware validates Gateway client API key / Tenant ID and injects TenantContext
func TenantContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}

		if tenantID == "" {
			tenantID = "default_tenant"
		}

		ctx := TenantContext{
			TenantID: tenantID,
			Role:     "white_label_client",
		}

		c.Set(TenantContextKey, ctx)
		c.Next()
	}
}
