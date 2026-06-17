package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/authorizerdev/authorizer/gen/openapi"
)

// spaBuildCacheMiddleware sets cache headers for SPA build assets:
//   - "index.js" / "main.css" (unhashed entry points the shell HTML loads
//     by name) → no-cache, so browsers always pick up new chunk references
//     after a deploy.
//   - everything else (content-hashed chunks, immutable assets) → long-lived
//     immutable cache, since a content change produces a new filename.
func spaBuildCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		base := path.Base(c.Request.URL.Path)
		if base == "index.js" || base == "main.css" {
			c.Header("Cache-Control", "no-cache, must-revalidate")
		} else {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	}
}

// NewRouter creates new gin router
func (s *server) NewRouter() *gin.Engine {
	router := gin.New()
	// Restrict the set of proxies whose forwarded headers are honoured.
	// When TrustedProxies is empty/nil, gin trusts NO proxies and falls back
	// to RemoteAddr — preventing X-Forwarded-For spoofing for rate limiting,
	// audit logs, and CSRF same-origin comparisons.
	var trustedProxies []string
	if s.Dependencies.AppConfig != nil {
		trustedProxies = s.Dependencies.AppConfig.TrustedProxies
	}
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		s.Dependencies.Log.Warn().Err(err).Msg("failed to apply trusted proxies; falling back to gin defaults")
	}
	router.Use(gin.Recovery())

	router.Use(s.Dependencies.HTTPProvider.SecurityHeadersMiddleware())
	router.Use(s.Dependencies.HTTPProvider.LoggerMiddleware())
	router.Use(s.Dependencies.HTTPProvider.MetricsMiddleware())
	router.Use(s.Dependencies.HTTPProvider.ContextMiddleware())
	router.Use(s.Dependencies.HTTPProvider.CORSMiddleware())
	router.Use(s.Dependencies.HTTPProvider.RateLimitMiddleware())
	router.Use(s.Dependencies.HTTPProvider.CSRFMiddleware())

	multiTenant := s.Dependencies.AppConfig != nil && s.Dependencies.AppConfig.EnableMultiTenant

	router.GET("/health", s.Dependencies.HTTPProvider.HealthHandler())
	router.GET("/healthz", s.Dependencies.HTTPProvider.HealthHandler())
	router.GET("/readyz", s.Dependencies.HTTPProvider.ReadyHandler())

	if multiTenant {
		tenant := router.Group("/:tenant_id")
		tenant.Use(s.Dependencies.HTTPProvider.TenantMiddleware())
		tenant.Use(s.Dependencies.HTTPProvider.ClientCheckMiddleware())
		s.mountAuthRoutes(tenant)
		tenant.POST("/_warm", s.Dependencies.HTTPProvider.WarmTenantHandler())
	} else {
		router.Use(s.Dependencies.HTTPProvider.ClientCheckMiddleware())
		s.mountAuthRoutes(router)
	}

	// Set up template functions for JSON encoding.
	// Escape </script> and <!-- to prevent script injection in <script> blocks.
	router.SetFuncMap(template.FuncMap{
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			s := string(a)
			s = strings.ReplaceAll(s, "</", `<\/`)
			s = strings.ReplaceAll(s, "<!--", `<\!--`)
			return template.JS(s)
		},
	})
	router.LoadHTMLGlob("web/templates/*")

	// SPA fallback: any unmatched GET inside /app/ or /dashboard/ serves the
	// SPA shell so deep links and browser refresh on multi-segment routes
	// (e.g. /dashboard/authorization/resources) don't return 404. Static
	// routes (/build, /favicon_io, /public) and the explicit /, /:page
	// handlers above take precedence; this only catches the multi-segment
	// gap. Non-GET methods and other paths fall through to gin's default
	// 404 handler.
	dashboardHandler := s.Dependencies.HTTPProvider.DashboardHandler()
	appHandler := s.Dependencies.HTTPProvider.AppHandler()
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.AbortWithStatus(404)
			return
		}
		path := c.Request.URL.Path
		switch {
		case strings.HasPrefix(path, "/dashboard/"):
			dashboardHandler(c)
		case strings.HasPrefix(path, "/app/"):
			appHandler(c)
		default:
			c.AbortWithStatus(404)
		}
	})
	return router
}

func (s *server) mountAuthRoutes(group gin.IRouter) {
	group.GET("/", s.Dependencies.HTTPProvider.RootHandler())
	group.POST("/graphql", s.Dependencies.HTTPProvider.GraphqlHandler())
	group.GET("/playground", s.Dependencies.HTTPProvider.PlaygroundHandler())
	group.GET("/oauth_login/:oauth_provider", s.Dependencies.HTTPProvider.OAuthLoginHandler())
	group.GET("/oauth_callback/:oauth_provider", s.Dependencies.HTTPProvider.OAuthCallbackHandler())
	group.POST("/oauth_callback/:oauth_provider", s.Dependencies.HTTPProvider.OAuthCallbackHandler())
	group.GET("/verify_email", s.Dependencies.HTTPProvider.VerifyEmailHandler())
	group.GET("/.well-known/openid-configuration", s.Dependencies.HTTPProvider.OpenIDConfigurationHandler())
	group.GET("/.well-known/jwks.json", s.Dependencies.HTTPProvider.JWKsHandler())
	group.GET("/authorize", s.Dependencies.HTTPProvider.AuthorizeHandler())
	group.GET("/userinfo", s.Dependencies.HTTPProvider.UserInfoHandler())
	group.GET("/logout", s.Dependencies.HTTPProvider.LogoutHandler())
	group.POST("/logout", s.Dependencies.HTTPProvider.LogoutHandler())
	group.POST("/oauth/token", s.Dependencies.HTTPProvider.TokenHandler())
	group.POST("/oauth/revoke", s.Dependencies.HTTPProvider.RevokeRefreshTokenHandler())
	group.POST("/oauth/introspect", s.Dependencies.HTTPProvider.IntrospectHandler())

	if s.gatewayHandler != nil {
		gw := gin.WrapH(s.gatewayHandler)
		group.Any("/v1/*path", gw)
		group.GET("/openapi.json", func(c *gin.Context) {
			c.Data(http.StatusOK, "application/json", openapi.Spec())
		})
	}

	app := group.Group("/app")
	{
		app.Static("/favicon_io", "web/app/favicon_io")
		appBuild := app.Group("/build")
		appBuild.Use(spaBuildCacheMiddleware())
		appBuild.Static("", "web/app/build")
		app.GET("/", s.Dependencies.HTTPProvider.AppHandler())
		app.GET("/:page", s.Dependencies.HTTPProvider.AppHandler())
	}

	dashboard := group.Group("/dashboard")
	{
		dashboard.Static("/favicon_io", "web/dashboard/favicon_io")
		dashboardBuild := dashboard.Group("/build")
		dashboardBuild.Use(spaBuildCacheMiddleware())
		dashboardBuild.Static("", "web/dashboard/build")
		dashboard.Static("/public", "web/dashboard/public")
		dashboard.GET("/", s.Dependencies.HTTPProvider.DashboardHandler())
		dashboard.GET("/:page", s.Dependencies.HTTPProvider.DashboardHandler())
	}
}
