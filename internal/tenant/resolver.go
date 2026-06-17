package tenant

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const PlatformTenantID = "platform"

// ResolveInput describes how to map tenants to database URLs and credentials.
type ResolveInput struct {
	PlatformDatabaseURL string
	PlatformClientID    string
	PlatformClientSecret string
	AppDatabaseURL      func(dbName string) string
}

// Resolver looks up tenant metadata from the platform database.
type Resolver struct {
	input ResolveInput
	db    *sql.DB
}

// NewResolver connects to the platform database for tenant lookups.
func NewResolver(input ResolveInput) (*Resolver, error) {
	db, err := sql.Open("pgx", input.PlatformDatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("platform database ping: %w", err)
	}
	return &Resolver{input: input, db: db}, nil
}

// Close closes the platform lookup connection.
func (r *Resolver) Close() error {
	if r.db == nil {
		return nil
	}
	return r.db.Close()
}

// Resolve returns tenant config for the given tenant id.
func (r *Resolver) Resolve(ctx context.Context, tenantID string) (Config, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return Config{}, fmt.Errorf("tenant id is required")
	}
	if tenantID == PlatformTenantID {
		return Config{
			TenantID:     PlatformTenantID,
			ClientID:     r.input.PlatformClientID,
			ClientSecret: r.input.PlatformClientSecret,
			DatabaseURL:  r.input.PlatformDatabaseURL,
		}, nil
	}
	var dbName, clientID, clientSecret sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT db_name, auth_client_id, auth_client_secret
		FROM apps WHERE id::text = $1 OR slug = $1
	`, tenantID).Scan(&dbName, &clientID, &clientSecret)
	if err != nil {
		return Config{}, fmt.Errorf("tenant %q not found: %w", tenantID, err)
	}
	if !dbName.Valid || dbName.String == "" {
		return Config{}, fmt.Errorf("tenant %q has no database", tenantID)
	}
	if r.input.AppDatabaseURL == nil {
		return Config{}, fmt.Errorf("app database resolver not configured")
	}
	cfg := Config{
		TenantID:    tenantID,
		DatabaseURL: r.input.AppDatabaseURL(dbName.String),
	}
	if clientID.Valid {
		cfg.ClientID = clientID.String
	}
	if clientSecret.Valid {
		cfg.ClientSecret = clientSecret.String
	}
	if cfg.ClientID == "" {
		cfg.ClientID = r.input.PlatformClientID
	}
	if cfg.ClientSecret == "" {
		cfg.ClientSecret = r.input.PlatformClientSecret
	}
	return cfg, nil
}
