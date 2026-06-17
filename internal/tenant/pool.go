package tenant

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/authorizerdev/authorizer/internal/config"
	"github.com/authorizerdev/authorizer/internal/storage"
	"github.com/authorizerdev/authorizer/internal/storage/db/sql"
)

type entry struct {
	provider storage.Provider
	cfg      Config
	lastUsed time.Time
	permanent bool
}

// Pool lazily creates and caches storage providers per tenant.
type Pool struct {
	baseConfig *config.Config
	log        *zerolog.Logger
	resolver   *Resolver
	ttl        time.Duration
	entries    sync.Map
}

// NewPool creates a tenant storage pool.
func NewPool(baseConfig *config.Config, log *zerolog.Logger, resolver *Resolver, ttl time.Duration) *Pool {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	p := &Pool{
		baseConfig: baseConfig,
		log:        log,
		resolver:   resolver,
		ttl:        ttl,
	}
	go p.pruneLoop()
	return p
}

// Get resolves tenant config and returns a storage provider.
func (p *Pool) Get(ctx context.Context, tenantID string) (storage.Provider, Config, error) {
	if tenantID == PlatformTenantID {
		if value, ok := p.entries.Load(PlatformTenantID); ok {
			e := value.(*entry)
			e.lastUsed = time.Now()
			return e.provider, e.cfg, nil
		}
	}
	if value, ok := p.entries.Load(tenantID); ok {
		e := value.(*entry)
		e.lastUsed = time.Now()
		return e.provider, e.cfg, nil
	}
	cfg, err := p.resolver.Resolve(ctx, tenantID)
	if err != nil {
		return nil, Config{}, err
	}
	provider, err := p.createProvider(cfg.DatabaseURL)
	if err != nil {
		return nil, Config{}, err
	}
	e := &entry{
		provider:  provider,
		cfg:       cfg,
		lastUsed:  time.Now(),
		permanent: tenantID == PlatformTenantID,
	}
	p.entries.Store(tenantID, e)
	if tenantID == PlatformTenantID {
		p.log.Info().Str("tenant", tenantID).Msg("initialized permanent platform storage pool")
	} else {
		p.log.Info().Str("tenant", tenantID).Str("db", cfg.DatabaseURL).Msg("initialized tenant storage pool")
	}
	return provider, cfg, nil
}

// Warm ensures a tenant pool entry exists (used on app provisioning).
func (p *Pool) Warm(ctx context.Context, tenantID string) error {
	_, _, err := p.Get(ctx, tenantID)
	return err
}

func (p *Pool) createProvider(databaseURL string) (storage.Provider, error) {
	cfg := *p.baseConfig
	cfg.DatabaseURL = databaseURL
	return sql.NewProvider(&cfg, &sql.Dependencies{Log: p.log})
}

// Close closes all pooled providers.
func (p *Pool) Close() error {
	var first error
	p.entries.Range(func(key, value any) bool {
		e := value.(*entry)
		if err := e.provider.Close(); err != nil && first == nil {
			first = err
		}
		p.entries.Delete(key)
		return true
	})
	return first
}

func (p *Pool) pruneLoop() {
	ticker := time.NewTicker(p.ttl / 2)
	defer ticker.Stop()
	for range ticker.C {
		p.prune()
	}
}

func (p *Pool) prune() {
	cutoff := time.Now().Add(-p.ttl)
	p.entries.Range(func(key, value any) bool {
		e := value.(*entry)
		if e.permanent || e.lastUsed.After(cutoff) {
			return true
		}
		_ = e.provider.Close()
		p.entries.Delete(key)
		p.log.Info().Str("tenant", fmt.Sprint(key)).Msg("pruned idle tenant storage pool")
		return true
	})
}

// PlatformHealthCheck verifies the platform tenant is reachable.
func (p *Pool) PlatformHealthCheck(ctx context.Context) error {
	provider, _, err := p.Get(ctx, PlatformTenantID)
	if err != nil {
		return err
	}
	return provider.HealthCheck(ctx)
}
