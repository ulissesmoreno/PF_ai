package postgres

import (
	"context"
	"database/sql"
	"time"

	"pf-ai/src/domain"
)

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
}

type Database struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Database {
	return Database{db: db}
}

func (d Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d Database) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

type ProviderRepository struct {
	store Store
}

func NewProviderRepository(store Store) ProviderRepository {
	return ProviderRepository{store: store}
}

func (r ProviderRepository) SaveProvider(ctx context.Context, provider domain.ModelProvider) error {
	_, err := r.store.ExecContext(ctx, `
INSERT INTO model_providers (id, name, mode, endpoint, model, secret_ref, local_runtime)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    mode = EXCLUDED.mode,
    endpoint = EXCLUDED.endpoint,
    model = EXCLUDED.model,
    secret_ref = EXCLUDED.secret_ref,
    local_runtime = EXCLUDED.local_runtime`,
		provider.ID,
		provider.Name,
		string(provider.Mode),
		provider.Endpoint,
		provider.Model,
		provider.SecretRef,
		provider.LocalRuntime,
	)
	return err
}

func (r ProviderRepository) ListProviders(ctx context.Context) ([]domain.ModelProvider, error) {
	rows, err := r.store.QueryContext(ctx, `
SELECT id, name, mode, endpoint, model, secret_ref, local_runtime FROM model_providers
ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := make([]domain.ModelProvider, 0)
	for rows.Next() {
		var provider domain.ModelProvider
		var mode string
		if err := rows.Scan(
			&provider.ID,
			&provider.Name,
			&mode,
			&provider.Endpoint,
			&provider.Model,
			&provider.SecretRef,
			&provider.LocalRuntime,
		); err != nil {
			return nil, err
		}
		provider.Mode = domain.ProviderMode(mode)
		providers = append(providers, provider)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return providers, nil
}

type AgentRepository struct {
	store Store
}

func NewAgentRepository(store Store) AgentRepository {
	return AgentRepository{store: store}
}

func (r AgentRepository) SaveAgent(ctx context.Context, agent domain.AgentDefinition) error {
	_, err := r.store.ExecContext(ctx, `
INSERT INTO agents (id, name, role, seniority, provider_id, description)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    role = EXCLUDED.role,
    seniority = EXCLUDED.seniority,
    provider_id = EXCLUDED.provider_id,
    description = EXCLUDED.description`,
		agent.ID,
		agent.Name,
		agent.Role,
		string(agent.Seniority),
		agent.ProviderID,
		agent.Description,
	)
	return err
}

func (r AgentRepository) ListAgents(ctx context.Context) ([]domain.AgentDefinition, error) {
	rows, err := r.store.QueryContext(ctx, `
SELECT id, name, role, seniority, provider_id, description FROM agents
ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := make([]domain.AgentDefinition, 0)
	for rows.Next() {
		var agent domain.AgentDefinition
		var seniority string
		if err := rows.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Role,
			&seniority,
			&agent.ProviderID,
			&agent.Description,
		); err != nil {
			return nil, err
		}
		agent.Seniority = domain.Seniority(seniority)
		agents = append(agents, agent)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return agents, nil
}

type SessionRepository struct {
	store Store
}

func NewSessionRepository(store Store) SessionRepository {
	return SessionRepository{store: store}
}

func (r SessionRepository) SaveSession(ctx context.Context, session domain.Session) error {
	_, err := r.store.ExecContext(ctx, `
INSERT INTO sessions (id, subject, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE SET
    subject = EXCLUDED.subject,
    expires_at = EXCLUDED.expires_at`,
		session.ID,
		session.Subject,
		session.ExpiresAt.UTC().Format(time.RFC3339),
	)
	return err
}

type AuditRepository struct {
	store Store
}

func NewAuditRepository(store Store) AuditRepository {
	return AuditRepository{store: store}
}

func (r AuditRepository) SaveAuditEvent(ctx context.Context, event domain.AuditEvent) error {
	_, err := r.store.ExecContext(ctx, `
INSERT INTO audit_events (actor, action, target, metadata)
VALUES ($1, $2, $3, $4::jsonb)`,
		event.Actor,
		event.Action,
		event.Target,
		event.MetadataJSON,
	)
	return err
}
