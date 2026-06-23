package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/authorizerdev/authorizer/internal/config"
	"github.com/authorizerdev/authorizer/internal/constants"
	"github.com/authorizerdev/authorizer/internal/http_handlers"
	"github.com/authorizerdev/authorizer/internal/storage"
	"github.com/authorizerdev/authorizer/internal/tenant"
)

func setupMultiTenantStorage(cfg *config.Config, log *zerolog.Logger) (storage.Provider, *tenant.Pool, *tenant.Resolver, error) {
	if strings.TrimSpace(cfg.PlatformDatabaseURL) == "" {
		return nil, nil, nil, fmt.Errorf("--platform-database-url is required in multi-tenant mode")
	}
	if cfg.DatabaseType == "" {
		cfg.DatabaseType = constants.DbTypePostgres
	}
	if cfg.DatabaseType != constants.DbTypePostgres {
		return nil, nil, nil, fmt.Errorf("multi-tenant mode currently supports postgres only")
	}
	appURL := func(dbName string) string {
		host := cfg.AppDatabaseHost
		port := cfg.AppDatabasePort
		user := cfg.AppDatabaseUser
		pass := cfg.AppDatabasePassword
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
	}
	resolver, err := tenant.NewResolver(tenant.ResolveInput{
		PlatformDatabaseURL:  cfg.PlatformDatabaseURL,
		PlatformClientID:     cfg.ClientID,
		PlatformClientSecret: cfg.ClientSecret,
		AppDatabaseURL:       appURL,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	ttl := time.Duration(cfg.PoolTTLMinutes) * time.Minute
	if ttl <= 0 {
		ttl = 120 * time.Minute
	}
	pool := tenant.NewPool(cfg, log, resolver, ttl)
	if err := pool.Warm(context.Background(), tenant.PlatformTenantID); err != nil {
		_ = pool.Close()
		_ = resolver.Close()
		return nil, nil, nil, fmt.Errorf("warm platform tenant: %w", err)
	}
	log.Info().Msg("multi-tenant storage pool initialized (platform warmed)")
	return http_handlers.ContextBoundStorage{}, pool, resolver, nil
}
