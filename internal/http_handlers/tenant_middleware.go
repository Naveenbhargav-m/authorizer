package http_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/authorizerdev/authorizer/internal/tenant"
	"github.com/authorizerdev/authorizer/internal/utils"
)

// TenantMiddleware resolves tenant storage from the URL path prefix /:tenant_id.
func (h *httpProvider) TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := strings.TrimSpace(c.Param("tenant_id"))
		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":             "invalid_tenant",
				"error_description": "tenant id is required in the URL path",
			})
			return
		}
		if h.TenantPool == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "tenant pool not configured",
			})
			return
		}
		provider, cfg, err := h.TenantPool.Get(c.Request.Context(), tenantID)
		if err != nil {
			h.Log.Error().Err(err).Str("tenant", tenantID).Msg("tenant resolution failed")
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":             "tenant_not_found",
				"error_description": err.Error(),
			})
			return
		}
		ctx := c.Request.Context()
		ctx = tenant.WithID(ctx, tenantID)
		ctx = tenant.WithConfig(ctx, cfg)
		ctx = tenant.WithStorage(ctx, provider)
		ctx = utils.ContextWithGin(ctx, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// WarmTenantHandler pre-initializes tenant storage (used after app provisioning).
func (h *httpProvider) WarmTenantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := strings.TrimSpace(c.Param("tenant_id"))
		if tenantID == "" || h.TenantPool == nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if c.GetHeader("X-Authorizer-Admin-Secret") != h.Config.AdminSecret {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err := h.TenantPool.Warm(c.Request.Context(), tenantID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "warmed", "tenant": tenantID})
	}
}
