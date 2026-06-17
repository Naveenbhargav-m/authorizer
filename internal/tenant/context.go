package tenant

import (
	"context"
	"fmt"

	"github.com/authorizerdev/authorizer/internal/storage"
)

type ctxKey int

const (
	storageKey ctxKey = iota
	configKey
	idKey
)

// Config holds per-tenant OAuth client credentials.
type Config struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	DatabaseURL  string
}

// WithStorage attaches a tenant storage provider to the context.
func WithStorage(ctx context.Context, p storage.Provider) context.Context {
	return context.WithValue(ctx, storageKey, p)
}

// WithConfig attaches tenant OAuth config to the context.
func WithConfig(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

// WithID attaches the tenant identifier to the context.
func WithID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, idKey, tenantID)
}

// IDFromContext returns the tenant id from context.
func IDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(idKey).(string)
	return v
}

// ConfigFromContext returns per-tenant OAuth config.
func ConfigFromContext(ctx context.Context) (Config, bool) {
	cfg, ok := ctx.Value(configKey).(Config)
	return cfg, ok
}

// StorageFromContext returns the tenant-bound storage provider.
func StorageFromContext(ctx context.Context) (storage.Provider, error) {
	p, ok := ctx.Value(storageKey).(storage.Provider)
	if !ok || p == nil {
		return nil, fmt.Errorf("tenant storage not found in context")
	}
	return p, nil
}
