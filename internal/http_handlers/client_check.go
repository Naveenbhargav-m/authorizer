package http_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/authorizerdev/authorizer/internal/metrics"
	"github.com/authorizerdev/authorizer/internal/tenant"
)

// ClientCheckMiddleware is a middleware to verify the client ID
// Note: client ID is passed in the header.
// An empty client ID is intentionally allowed for routes that don't require it
// (e.g., OAuth callbacks, JWKS, OpenID configuration, health checks).
// The middleware only rejects requests with an explicitly wrong client ID.
func (h *httpProvider) ClientCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.Log.With().Str("func", "ClientCheckMiddleware").
			Str("path", c.Request.URL.Path).
			Logger()
		// Only check client ID for graphql route [Most relevant route for client ID check]
		if c.Request.URL.Path != "/graphql" && !strings.Contains(c.Request.URL.Path, "/graphql") {
			c.Next()
			return
		}
		clientID := c.Request.Header.Get("X-Authorizer-Client-ID")
		// Allowing requests without client ID header for backward compatibility.
		// The dashboard and other first-party clients may not always send this
		// header, so an empty value must pass through to the GraphQL handler.
		if clientID == "" {
			log.Debug().Msg("request received without client ID header")
			metrics.RecordClientIDHeaderMissing()
			c.Next()
			return
		}

		expectedClientID := h.Config.ClientID
		if tenantCfg, ok := tenant.ConfigFromContext(c.Request.Context()); ok && tenantCfg.ClientID != "" {
			expectedClientID = tenantCfg.ClientID
		}

		if clientID != expectedClientID {
			// Record metric for client-id mismatch, but skip dashboard and app UI routes
			// as those are internal requests that should not trigger security alerts.
			metrics.RecordSecurityEvent("client_id_mismatch", "invalid_client_id")
			log.Debug().Str("client_id", clientID).Msg("Client ID is invalid")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":             "invalid_client_id",
				"error_description": "The client id is invalid",
			})
			return
		}

		c.Next()
	}
}
