package infrastructure_test

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"pf-ai/src/domain"
	"pf-ai/src/infrastructure/postgres"
)

type recordedSQL struct {
	query string
	args  []any
	rows  postgres.Rows
	err   error
}

func (r *recordedSQL) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	r.query = query
	r.args = args
	return nil, r.err
}

func (r *recordedSQL) QueryContext(_ context.Context, query string, args ...any) (postgres.Rows, error) {
	r.query = query
	r.args = args
	return r.rows, r.err
}

type fakeRows struct {
	values [][]any
	index  int
	err    error
}

func (r *fakeRows) Next() bool {
	return r.index < len(r.values)
}

func (r *fakeRows) Scan(dest ...any) error {
	if r.index >= len(r.values) {
		return errors.New("scan past end")
	}
	row := r.values[r.index]
	r.index++
	for i := range dest {
		switch d := dest[i].(type) {
		case *string:
			*d = row[i].(string)
		default:
			return errors.New("unsupported scan destination")
		}
	}
	return nil
}

func (r *fakeRows) Close() error {
	return nil
}

func (r *fakeRows) Err() error {
	return r.err
}

func TestPostgresProviderRepositorySavesWithUpsert(t *testing.T) {
	db := &recordedSQL{}
	repo := postgres.NewProviderRepository(db)
	provider, err := domain.NewModelProvider("api", "API", domain.ProviderModeAPI, "https://api.example.test", "gpt", "PF_AI_API_KEY", "")
	if err != nil {
		t.Fatalf("unexpected provider error: %v", err)
	}

	if err := repo.SaveProvider(context.Background(), provider); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	if !strings.Contains(db.query, "INSERT INTO model_providers") {
		t.Fatalf("expected model_providers insert, got %s", db.query)
	}
	if !strings.Contains(db.query, "ON CONFLICT (id) DO UPDATE") {
		t.Fatalf("expected upsert query, got %s", db.query)
	}
	if got := db.args[5]; got != "PF_AI_API_KEY" {
		t.Fatalf("expected secret ref argument, got %v", got)
	}
}

func TestPostgresAgentRepositoryListsAgents(t *testing.T) {
	db := &recordedSQL{
		rows: &fakeRows{values: [][]any{
			{"ceo", "CEO", "orchestration", "Senior", "api", "orchestrates"},
		}},
	}
	repo := postgres.NewAgentRepository(db)

	agents, err := repo.ListAgents(context.Background())
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}

	if len(agents) != 1 {
		t.Fatalf("expected one agent, got %d", len(agents))
	}
	if agents[0].Seniority != domain.SenioritySenior {
		t.Fatalf("expected Senior, got %s", agents[0].Seniority)
	}
	if !strings.Contains(db.query, "SELECT id, name, role, seniority, provider_id, description FROM agents") {
		t.Fatalf("expected agents select, got %s", db.query)
	}
}

func TestPostgresSessionRepositorySavesSession(t *testing.T) {
	db := &recordedSQL{}
	repo := postgres.NewSessionRepository(db)
	expiresAt := time.Date(2026, 5, 5, 19, 0, 0, 0, time.UTC)
	session, err := domain.NewSession("session-1", "ulisses", expiresAt)
	if err != nil {
		t.Fatalf("unexpected session error: %v", err)
	}

	if err := repo.SaveSession(context.Background(), session); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	if !strings.Contains(db.query, "INSERT INTO sessions") {
		t.Fatalf("expected sessions insert, got %s", db.query)
	}
	if got := db.args[2]; got != "2026-05-05T19:00:00Z" {
		t.Fatalf("expected UTC expiration, got %v", got)
	}
}

func TestPostgresAuditRepositorySavesAuditEvent(t *testing.T) {
	db := &recordedSQL{}
	repo := postgres.NewAuditRepository(db)
	event, err := domain.NewAuditEvent("ceo", "handoff.created", "handoff-1", `{"phase":"1"}`)
	if err != nil {
		t.Fatalf("unexpected audit event error: %v", err)
	}

	if err := repo.SaveAuditEvent(context.Background(), event); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	if !strings.Contains(db.query, "INSERT INTO audit_events") {
		t.Fatalf("expected audit_events insert, got %s", db.query)
	}
	if got := db.args[3]; got != `{"phase":"1"}` {
		t.Fatalf("expected metadata json, got %v", got)
	}
}
