package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"pf_ai/context_store"
)

func TestCreateCEOKickoffWithoutProjectCreatesOnboardingHandoff(t *testing.T) {
	workspace := t.TempDir()
	store := openMainTestStore(t, workspace)
	defer store.Close()

	cfg := Config{
		WorkspaceRoot:   workspace,
		AgentHandoffDir: filepath.Join(workspace, ".agent_handoff"),
	}
	if err := createCEOKickoff(cfg, store); err != nil {
		t.Fatalf("createCEOKickoff: %v", err)
	}

	handoff := readOnlyHandoff(t, cfg.AgentHandoffDir)
	if handoff.Header.Intent != "PROJECT_ONBOARDING" {
		t.Fatalf("intent = %q, want PROJECT_ONBOARDING", handoff.Header.Intent)
	}
	if handoff.Header.TaskRef != "ONBOARDING-1" {
		t.Fatalf("task_ref = %q, want ONBOARDING-1", handoff.Header.TaskRef)
	}
	if handoff.Payload.Action != "ask_human" || len(handoff.Payload.Questions) == 0 {
		t.Fatalf("payload = %#v", handoff.Payload)
	}
}

func TestCreateCEOKickoffWithProjectCreatesPhaseKickoff(t *testing.T) {
	workspace := t.TempDir()
	store := openMainTestStore(t, workspace)
	defer store.Close()
	if _, err := store.CreateProject(context_store.Project{Name: "PF AI", Slug: "pf-ai"}); err != nil {
		t.Fatalf("CreateProject: %v", err)
	}

	cfg := Config{
		WorkspaceRoot:   workspace,
		AgentHandoffDir: filepath.Join(workspace, ".agent_handoff"),
	}
	if err := createCEOKickoff(cfg, store); err != nil {
		t.Fatalf("createCEOKickoff: %v", err)
	}

	handoff := readOnlyHandoff(t, cfg.AgentHandoffDir)
	if handoff.Header.Intent != "PHASE_KICKOFF" {
		t.Fatalf("intent = %q, want PHASE_KICKOFF", handoff.Header.Intent)
	}
	if handoff.Header.TaskRef != "PHASE-1" {
		t.Fatalf("task_ref = %q, want PHASE-1", handoff.Header.TaskRef)
	}
}

func openMainTestStore(t *testing.T, workspace string) *context_store.Store {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	migrationsDir := filepath.Clean(filepath.Join(wd, "..", "db", "migrations"))
	store, err := context_store.Open(filepath.Join(workspace, "data", "test.db"), migrationsDir, workspace)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	return store
}

func readOnlyHandoff(t *testing.T, dir string) struct {
	Header struct {
		TaskRef string `json:"task_ref"`
		Intent  string `json:"intent"`
	} `json:"header"`
	Payload ceoKickoffPayload `json:"payload"`
} {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("handoffs = %d, want 1", len(entries))
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	var handoff struct {
		Header struct {
			TaskRef string `json:"task_ref"`
			Intent  string `json:"intent"`
		} `json:"header"`
		Payload ceoKickoffPayload `json:"payload"`
	}
	if err := json.Unmarshal(data, &handoff); err != nil {
		t.Fatal(err)
	}
	return handoff
}
